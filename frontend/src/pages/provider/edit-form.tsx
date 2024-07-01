import {
  ProviderItem,
  useProviderEditMutation,
} from '@/services/provider.service';
import {
  Form,
  FormInstance,
  Input,
  Message,
  Modal,
  Select,
} from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';
import { FunctionComponent, useEffect, useRef, useState } from 'react';
import { platformSupportTag } from '.';
import Tencent from './platform/tencent';
import Cloudflare from './platform/cloudflare';
import Qiniu from './platform/qiniu';
import Aliyun from './platform/aliyun';
import Kubernetes from './platform/kubernetes';

export const platformForm: {
  [key: string]: {
    name: string;
    component: FunctionComponent<{ update: boolean }>;
  };
} = {
  tencent: {
    name: '腾讯云',
    component: Tencent,
  },
  aliyun: {
    name: '阿里云',
    component: Aliyun,
  },
  cloudflare: {
    name: 'Cloudflare',
    component: Cloudflare,
  },
  qiniu: {
    name: '七牛云',
    component: Qiniu,
  },
  kubernetes: {
    name: 'Kubernetes',
    component: Kubernetes,
  },
};

const EditForm: React.FC<{
  visible: boolean;
  onOk: () => void;
  onCancel: () => void;
  data?: ProviderItem;
}> = (props) => {
  const [form] = Form.useForm();
  const [platform, setPlatform] = useState<string | undefined>(
    props.data?.platform,
  );
  const [editProvider, { isLoading }] = useProviderEditMutation();
  useEffect(() => {
    if (props.data) {
      form.setFieldsValue(props.data);
      setPlatform(props.data.platform);
    } else {
      form.resetFields();
      setPlatform(undefined);
    }
  }, [props.data]);

  const PlaformForm = platformForm[platform];

  return (
    <Modal
      title={props.data ? '编辑' + props.data.name : '新增厂商'}
      visible={props.visible}
      style={{ width: 600 }}
      confirmLoading={isLoading}
      onOk={async () => {
        form.validate().then((res) => {
          const values = form.getFieldsValue();
          const secret = {};
          Object.keys(values['secret']).forEach((key) => {
            if (values['secret'][key]) {
              secret[key] = values['secret'][key];
            }
          });
          editProvider({
            id: props.data?.id,
            name: values['name'],
            platform: values['platform'],
            secret: secret,
          })
            .unwrap()
            .then((res) => {
              Message.success('操作成功');
              props.onOk();
            });
        });
      }}
      onCancel={() => props.onCancel()}
    >
      <Form form={form} autoComplete="off">
        <FormItem field="name" label="名称" rules={[{ required: true }]}>
          <Input placeholder="请输入厂商名称" />
        </FormItem>
        <FormItem
          field="platform"
          label="平台"
          rules={[{ required: true }]}
          disabled={!!props.data?.id}
        >
          <Select
            placeholder="请选择厂商平台"
            value={platform}
            onChange={(val) => {
              setPlatform(val);
            }}
          >
            {Object.keys(platformForm).map((key) => (
              <Select.Option value={key} key={key}>
                {platformForm[key].name}
              </Select.Option>
            ))}
          </Select>
        </FormItem>
        <FormItem label="支持">{platformSupportTag(platform)}</FormItem>
        {(PlaformForm && (
          <PlaformForm.component update={props.data?.id ? true : false} />
        )) || <></>}
      </Form>
    </Modal>
  );
};

export default EditForm;

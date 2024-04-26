import { ProviderItem, useEditMutation } from '@/services/provider.service';
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

const platformForm: { [key: string]: FunctionComponent<{ update: boolean }> } =
  {
    tencent: Tencent,
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
  const [editProvider, { isLoading }] = useEditMutation();
  useEffect(() => {
    if (props.data) {
      form.setFieldsValue(props.data);
      setPlatform(props.data.platform);
    } else {
      form.resetFields();
      setPlatform(undefined);
    }
  }, [props.data]);

  const PlaformFormComponent = platformForm[platform];

  return (
    <Modal
      title={props.data ? '编辑' + props.data.name : '新增厂商'}
      visible={props.visible}
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
            <Select.Option value="tencent">腾讯云</Select.Option>
            <Select.Option value="cloudflare">Cloudflare</Select.Option>
          </Select>
        </FormItem>
        <FormItem label="支持">{platformSupportTag(platform)}</FormItem>
        {(PlaformFormComponent && (
          <PlaformFormComponent update={props.data?.id ? true : false} />
        )) || <></>}
      </Form>
    </Modal>
  );
};

export default EditForm;

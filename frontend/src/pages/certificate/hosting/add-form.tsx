import Kubernetes from '@/pages/certificate/hosting/platform/kubernetes';
import { useCdnListQuery } from '@/services/cdn.service';
import {
  useCertHostingAddMutation,
  useCertHostingQueryQuery,
} from '@/services/cert.service';
import { useProviderListQuery } from '@/services/provider.service';
import {
  Checkbox,
  Form,
  FormInstance,
  Input,
  InputTag,
  Message,
  Modal,
  Radio,
  Select,
} from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';
import { FunctionComponent, useEffect, useRef, useState } from 'react';

export const platformForm: {
  [key: string]: {
    name: string;
    component: FunctionComponent;
  };
} = {
  kubernetes: {
    name: 'Kubernetes',
    component: Kubernetes,
  },
};

const AddForm: React.FC<{
  visible: boolean;
  onOk: () => void;
  onCancel: () => void;
}> = (props) => {
  const [form] = Form.useForm();
  const [add, { isLoading }] = useCertHostingAddMutation();
  const { data: cdnData, isLoading: cdnIsLoading } = useCertHostingQueryQuery();
  const { data: providerData, isLoading: providerIsLoading } =
    useProviderListQuery();
  const [platform, setPlatform] = useState<string | undefined>();
  const [type, setType] = useState<string>('cdn');
  useEffect(() => {
    if (props.visible) {
      form.resetFields();
    }
  }, [props.visible]);

  const PlaformForm = platformForm[platform];

  return (
    <Modal
      title={'新增托管'}
      visible={props.visible}
      confirmLoading={isLoading}
      onOk={async () => {
        form.validate().then((res) => {
          const values = form.getFieldsValue();
          add({
            email: values['email'],
            type: values['type'] == 'cdn' ? 1 : 2,
            cdn_id: values['cdn_id'],
            provider_id: values['provider_id'],
            config: values['config'],
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
        <FormItem field="email" label="邮箱" rules={[{ required: true }]}>
          <Input placeholder="请输入邮箱" />
        </FormItem>
        <FormItem
          field="type"
          label="类型"
          rules={[{ required: true }]}
          defaultValue={type}
        >
          <Radio.Group
            value={type}
            type="button"
            defaultValue={type}
            onChange={(val) => {
              setType(val);
            }}
          >
            <Radio value="cdn">CDN</Radio>
            <Radio value="provider">厂商</Radio>
          </Radio.Group>
        </FormItem>
        {type == 'cdn' && (
          <FormItem field="cdn_id" label="CDN" rules={[{ required: true }]}>
            <Select loading={cdnIsLoading}>
              {cdnData?.data?.list?.map((item) => (
                <Select.Option
                  key={item.id}
                  value={item.id}
                  disabled={item.is_managed}
                >
                  {item.domain}
                </Select.Option>
              ))}
            </Select>
          </FormItem>
        )}
        {type == 'provider' && (
          <>
            <FormItem
              field="provider_id"
              label="厂商"
              rules={[{ required: true }]}
            >
              <Select
                loading={providerIsLoading}
                onChange={(v) => {
                  console.log(
                    providerData?.data?.list.find((item) => item.id == v),
                  );
                  setPlatform(
                    providerData?.data?.list.find((item) => item.id == v)
                      ?.platform,
                  );
                }}
              >
                {providerData?.data?.list?.map((item) => (
                  <Select.Option
                    key={item.id}
                    value={item.id}
                    disabled={!['kubernetes'].includes(item.platform)}
                  >
                    {item.name}
                  </Select.Option>
                ))}
              </Select>
            </FormItem>
            {PlaformForm && <PlaformForm.component />}
          </>
        )}
      </Form>
    </Modal>
  );
};

export default AddForm;

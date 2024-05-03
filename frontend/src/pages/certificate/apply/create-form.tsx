import { useCertCreateMutation } from '@/services/cert.service';
import {
  ProviderItem,
  useProviderEditMutation,
} from '@/services/provider.service';
import {
  Form,
  FormInstance,
  Input,
  InputTag,
  Message,
  Modal,
  Select,
} from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';
import { FunctionComponent, useEffect, useRef, useState } from 'react';

const CreateForm: React.FC<{
  visible: boolean;
  onOk: () => void;
  onCancel: () => void;
}> = (props) => {
  const [form] = Form.useForm();
  const [createCert, { isLoading }] = useCertCreateMutation();
  useEffect(() => {
    if (props.visible) {
      form.resetFields();
    }
  }, [props.visible]);

  return (
    <Modal
      title={'申请证书'}
      visible={props.visible}
      confirmLoading={isLoading}
      onOk={async () => {
        form.validate().then((res) => {
          const values = form.getFieldsValue();
          createCert({
            email: values['email'],
            domains: values['domains'],
          })
            .unwrap()
            .then((res) => {
              Message.success('申请请求创建成功，请等待申请');
              props.onOk();
            });
        });
      }}
      onCancel={() => props.onCancel()}
    >
      <Form form={form} autoComplete="off">
        <FormItem field="email" label="邮箱" rules={[{ required: true }]}>
          <Input placeholder="请输入申请邮箱" />
        </FormItem>
        <FormItem field="domains" label="域名" rules={[{ required: true }]}>
          <InputTag placeholder="请输入要申请的域名，按下回车添加" />
        </FormItem>
      </Form>
    </Modal>
  );
};

export default CreateForm;

import { Input } from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';

const Tencent: React.FC<{ update: boolean }> = ({ update }) => {
  return (
    <>
      <FormItem
        field="secret.secret_id"
        label="密钥ID"
        rules={update ? [] : [{ required: true }]}
      >
        <Input />
      </FormItem>
      <FormItem
        field="secret.secret_key"
        label="密钥"
        rules={update ? [] : [{ required: true }]}
      >
        <Input />
      </FormItem>
    </>
  );
};

export default Tencent;

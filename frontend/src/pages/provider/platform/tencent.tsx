import { Input, Space, Tag } from '@arco-design/web-react';
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

export const PlatformSupportTag = () => (
  <Space>
    <Tag color="green">CDN</Tag>
    <Tag color="arcoblue">DNS</Tag>
  </Space>
);

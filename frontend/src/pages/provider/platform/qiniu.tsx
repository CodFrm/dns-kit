import { Input, Space, Tag } from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';

const Qiniu: React.FC<{ update: boolean }> = ({ update }) => {
  return (
    <>
      <FormItem
        field="secret.access_key"
        label="AccessKey"
        rules={update ? [] : [{ required: true }]}
      >
        <Input />
      </FormItem>
      <FormItem
        field="secret.secret_key"
        label="SecretKey"
        rules={update ? [] : [{ required: true }]}
      >
        <Input />
      </FormItem>
    </>
  );
};

export default Qiniu;

export const PlatformSupportTag = () => (
  <Space>
    <Tag color="green">CDN</Tag>
  </Space>
);

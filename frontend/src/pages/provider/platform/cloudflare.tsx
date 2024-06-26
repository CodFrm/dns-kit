import { Input, Space, Tag } from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';

const Cloudflare: React.FC<{ update: boolean }> = ({ update }) => {
  return (
    <>
      <FormItem
        field="secret.token"
        label="Token"
        rules={update ? [] : [{ required: true }]}
      >
        <Input />
      </FormItem>
    </>
  );
};

export default Cloudflare;

export const PlatformSupportTag = () => (
  <Space>
   <Tag color="arcoblue">DNS</Tag>
  </Space>
);

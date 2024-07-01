import { Input, Space, Tag } from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';
import React from 'react';

const Kubernetes: React.FC<{ update: boolean }> = ({ update }) => {
  return (
    <>
      <FormItem field="secret.kube_config" label="KubeConfig">
        <Input.TextArea placeholder="为空将会读取默认配置" />
      </FormItem>
    </>
  );
};

export default Kubernetes;

export const PlatformSupportTag = () => (
  <Space>
    <Tag color="green">CDN</Tag>
  </Space>
);

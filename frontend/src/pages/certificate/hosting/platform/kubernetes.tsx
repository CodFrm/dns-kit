import { Input, Space, Tag } from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';
import React from 'react';

const Kubernetes: React.FC = () => {
  return (
    <>
      <FormItem
        field="config.namespace"
        label="命名空间"
        rules={[{ required: true }]}
      >
        <Input placeholder="secret部署的命名空间" />
      </FormItem>
      <FormItem field="config.domain" label="域名" rules={[{ required: true }]}>
        <Input placeholder="申请的域名" />
      </FormItem>
    </>
  );
};

export default Kubernetes;

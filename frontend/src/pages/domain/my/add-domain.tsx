import {
  useDomainAddMutation,
  useDomainQueryQuery,
} from '@/services/domain.service';
import { ProviderItem } from '@/services/provider.service';
import {
  Button,
  Form,
  Input,
  List,
  Message,
  Modal,
  Select,
  Space,
  Tag,
  Typography,
} from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';
import { useEffect, useState } from 'react';

const AddDomain: React.FC<{
  visible: boolean;
  onOk: () => void;
  onCancel: () => void;
  data?: ProviderItem;
}> = (props) => {
  const { data, isLoading } = useDomainQueryQuery();
  const [add, { isLoading: addIsLoading }] = useDomainAddMutation();
  const loading = isLoading || addIsLoading;

  return (
    <Modal
      title={'纳管域名'}
      visible={props.visible}
      confirmLoading={loading}
      onOk={async () => {
        props.onOk();
      }}
      onCancel={() => props.onCancel()}
    >
      <List
        style={{ width: '100%' }}
        size="small"
        header="域名列表"
        dataSource={data?.data.items}
        render={(item, index) => (
          <List.Item key={index}>
            <div className="flex flex-row justify-between">
              <Typography.Text>
                <Space>
                  {item.domain}
                  <Tag>{item.provider_name}</Tag>
                </Space>
              </Typography.Text>
              <Button
                type="primary"
                disabled={item.is_managed}
                onClick={() => {
                  add({
                    provider_id: item.provider_id,
                    domain_id: item.domain_id,
                    domain: item.domain,
                  });
                }}
                size="mini"
              >
                纳管
              </Button>
            </div>
          </List.Item>
        )}
      />
    </Modal>
  );
};

export default AddDomain;

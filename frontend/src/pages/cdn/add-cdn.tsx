import { useCdnAddMutation, useCdnQueryQuery } from '@/services/cdn.service';
import {
  Button,
  List,
  Message,
  Modal,
  Space,
  Tag,
  Typography,
} from '@arco-design/web-react';

const AddCDN: React.FC<{
  visible: boolean;
  onOk: () => void;
  onCancel: () => void;
}> = (props) => {
  const { data, isLoading } = useCdnQueryQuery();
  const [add, { isLoading: addIsLoading }] = useCdnAddMutation();
  const loading = isLoading || addIsLoading;

  return (
    <Modal
      title={'纳管CDN'}
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
        header="CDN列表"
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
                    id: item.id,
                    domain: item.domain,
                  })
                    .unwrap()
                    .then(() => {
                      Message.success('纳管成功');
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

export default AddCDN;

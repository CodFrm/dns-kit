import { useState } from 'react';
import {
  Button,
  Card,
  Message,
  Popconfirm,
  Space,
  Table,
} from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table';
import {
  IconDelete,
  IconPlus,
} from '@arco-design/web-react/icon';
import Title from '@arco-design/web-react/es/Typography/title';
import AddCDN from './add-cdn';
import { PiCertificateBold } from 'react-icons/pi';
import {
  CDNItem,
  useCdnDeleteMutation,
  useCdnListQuery,
} from '@/services/cdn.service';

function CDN() {
  const [visible, setVisible] = useState(false);
  const [addDomainVisible, setAddDomainVisible] = useState(false);
  const [editData, setEditData] = useState<CDNItem>();
  const { data, isLoading: listLoading } = useCdnListQuery();
  const [deleteCdn, { isLoading: deleteLoading }] = useCdnDeleteMutation();

  const isLoading = listLoading || deleteLoading;

  const columns: ColumnProps<CDNItem>[] = [
    {
      key: 'id',
      title: 'ID',
      dataIndex: 'id',
    },
    {
      key: 'provider_name',
      title: '供应商',
      dataIndex: 'provider_name',
    },
    { key: 'domain', title: '域名', dataIndex: 'domain' },
    {
      key: 'action',
      title: '操作',
      render(col, item) {
        return (
          <Space key={item.id}>
            <Button
              type="text"
              style={{ color: 'var(--color-text-2)' }}
              iconOnly
              icon={<PiCertificateBold />}
              onClick={() => {
                setEditData(item);
                setVisible(true);
              }}
            />
            <Popconfirm
              focusLock
              title="确定"
              content="确认删除吗？删除后相关的资源也会被删除"
              onOk={() => {
                deleteCdn(item.id)
                  .unwrap()
                  .then(() => {
                    Message.success('删除成功');
                  });
              }}
            >
              <Button
                type="text"
                style={{ color: 'var(--color-text-2)' }}
                iconOnly
                icon={<IconDelete />}
              />
            </Popconfirm>
          </Space>
        );
      },
    },
  ];

  return (
    <Card style={{ height: '80vh' }}>
      <Title heading={6}>我的域名</Title>
      <div className="flex flex-col">
        <div className="text-right">
          <AddCDN
            visible={addDomainVisible}
            onOk={() => {
              setAddDomainVisible(false);
            }}
            onCancel={() => {
              setAddDomainVisible(false);
            }}
          />
          <Button
            style={{ marginBottom: 10 }}
            type="primary"
            icon={<IconPlus />}
            onClick={() => {
              setAddDomainVisible(true);
            }}
          >
            纳管
          </Button>
        </div>
        <Table
          columns={columns}
          loading={isLoading}
          data={data?.data.list}
          border={{}}
          pagination={{
            pageSize: 20,
            total: data?.data?.total,
          }}
        />
      </div>
    </Card>
  );
}

export default CDN;

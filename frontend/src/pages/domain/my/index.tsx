import { useState } from 'react';
import { Button, Card, Message, Popconfirm, Space, Table } from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table';
import { IconDelete, IconDesktop, IconPlus } from '@arco-design/web-react/icon';
import Title from '@arco-design/web-react/es/Typography/title';
import RecordManager from './record-manager';
import AddDomain from './add-domain';
import {
  DomainItem,
  useDomainDeleteMutation,
  useDomainListQuery,
} from '@/services/domain.service';

function My() {
  const [visible, setVisible] = useState(false);
  const [addDomainVisible, setAddDomainVisible] = useState(false);
  const [editData, setEditData] = useState<DomainItem>();
  const { data, isLoading: listLoading } = useDomainListQuery();
  const [deleteDomain, { isLoading: deleteLoading }] =
    useDomainDeleteMutation();

  const isLoading = listLoading || deleteLoading;

  const columns: ColumnProps<DomainItem>[] = [
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
              icon={<IconDesktop />}
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
                deleteDomain(item.id)
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
          <AddDomain
            visible={addDomainVisible}
            onOk={() => {
              setAddDomainVisible(false);
            }}
            onCancel={() => {
              setAddDomainVisible(false);
            }}
          />
          <RecordManager
            visible={visible}
            onOk={() => {
              setVisible(false);
            }}
            onCancel={() => {
              setVisible(false);
            }}
            domain={editData}
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

export default My;

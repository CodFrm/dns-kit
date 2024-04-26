import { useState } from 'react';
import { Button, Card, Popconfirm, Space, Table } from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table';
import {
  IconDelete,
  IconDesktop,
  IconEdit,
  IconList,
  IconPlus,
} from '@arco-design/web-react/icon';
import Title from '@arco-design/web-react/es/Typography/title';
import EditForm from '@/pages/provider/edit-form';
import { ProviderItem } from '@/services/provider.service';
import RecordManager from './record-manager';

function My() {
  const [visible, setVisible] = useState(false);
  const [editData, setEditData] = useState<ProviderItem>();

  const columns: ColumnProps<ProviderItem>[] = [
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
    { key: 'name', title: '名称', dataIndex: 'name' },
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
              onOk={() => {}}
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
          <RecordManager
            visible={visible}
            onOk={() => {
              setVisible(false);
            }}
            onCancel={() => {
              setVisible(false);
            }}
            data={editData}
          />
          <Button
            style={{ marginBottom: 10 }}
            type="primary"
            icon={<IconPlus />}
            onClick={() => {
              setEditData(null);
              setVisible(true);
            }}
          >
            纳管
          </Button>
        </div>
        <Table
          columns={columns}
          // loading={isLoading}
          data={[{ id: 1, name: '1', platform: '1' }]}
          border={{}}
          pagination={{
            pageSize: 20,
            // total: data?.data?.total,
          }}
        />
      </div>
    </Card>
  );
}

export default My;

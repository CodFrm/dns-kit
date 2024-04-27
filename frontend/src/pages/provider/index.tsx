import { useState } from 'react';
import {
  Button,
  Card,
  Message,
  Popconfirm,
  Space,
  Table,
  Tag,
} from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table';
import { IconDelete, IconEdit, IconPlus } from '@arco-design/web-react/icon';
import Title from '@arco-design/web-react/es/Typography/title';
import EditForm from '@/pages/provider/edit-form';
import {
  ProviderItem,
  useProviderDeleteMutation,
  useProviderListQuery,
} from '@/services/provider.service';

export function platformSupportTag(platform: string) {
  switch (platform) {
    case 'tencent':
      return (
        <Space>
          <Tag color="green">CDN</Tag>
          <Tag color="arcoblue">DNS</Tag>
        </Space>
      );
  }
  return <></>;
}

function Provider() {
  const [visible, setVisible] = useState(false);
  const [editData, setEditData] = useState<ProviderItem>();
  const { data, isLoading } = useProviderListQuery();
  const [deleteProvider, {}] = useProviderDeleteMutation();
  const columns: ColumnProps<ProviderItem>[] = [
    {
      key: 'id',
      title: 'ID',
      dataIndex: 'id',
    },
    {
      key: 'name',
      title: '名称',
      dataIndex: 'name',
    },
    {
      key: 'platform',
      title: '平台',
      dataIndex: 'platform',
      render(col) {
        return { tencent: '腾讯云' }[col];
      },
    },
    {
      key: 'support',
      title: '支持',
      dataIndex: 'platform',
      render(col) {
        return platformSupportTag(col);
      },
    },
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
              icon={<IconEdit />}
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
                deleteProvider(item.id)
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
      <Title heading={6}>厂商管理</Title>
      <div className="flex flex-col">
        <div className="text-right">
          <EditForm
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
            添加
          </Button>
        </div>
        <Table
          columns={columns}
          loading={isLoading}
          data={data?.data?.list}
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

export default Provider;

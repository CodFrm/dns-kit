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
import {
  CertHostingItem,
  useCertHostigListQuery,
  useCertHostingDeleteMutation,
} from '@/services/cert.service';
import AddForm from './add-form';

function Hosting() {
  const [visible, setVisible] = useState(false);
  const [editData, setEditData] = useState<CertHostingItem>();
  const { data, isLoading } = useCertHostigListQuery();
  const [deleteProvider, {}] = useCertHostingDeleteMutation();
  const columns: ColumnProps<CertHostingItem>[] = [
    {
      key: 'id',
      title: 'ID',
      dataIndex: 'id',
    },
    {
      key: 'cdn',
      title: 'CDN',
      dataIndex: 'cdn',
    },
    {
      key: 'cert',
      title: '关联证书',
      dataIndex: 'cert_id',
    },
    {
      key: 'status',
      title: '状态',
      dataIndex: 'status',
      render(col) {
        switch (col) {
          case 1:
            return <Tag color="green">成功</Tag>;
          case 3:
            return <Tag color="blue">部署中</Tag>;
          case 4:
            return <Tag color="orangered">部署失败(等待重试)</Tag>;
          case 5:
            return <Tag color="red">部署失败</Tag>;
          default:
            return '未知';
        }
      },
    },
    {
      key: 'action',
      title: '操作',
      render(col, item) {
        return (
          <Space key={item.id}>
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
      <Title heading={6}>证书托管</Title>
      <div className="flex flex-col">
        <div className="text-right">
          <AddForm
            visible={visible}
            onOk={() => {
              setVisible(false);
            }}
            onCancel={() => {
              setVisible(false);
            }}
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

export default Hosting;

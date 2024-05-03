import { useState } from 'react';
import {
  Button,
  Card,
  Popconfirm,
  Space,
  Table,
  Tag,
} from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table';
import {
  IconDelete,
  IconDownload,
  IconPlus,
} from '@arco-design/web-react/icon';
import Title from '@arco-design/web-react/es/Typography/title';
import { CertItem, useCertListQuery } from '@/services/cert.service';
import CreateForm from './create-form';

function Apply() {
  const [visible, setVisible] = useState(false);
  const { data, isLoading } = useCertListQuery();
  const columns: ColumnProps<CertItem>[] = [
    {
      key: 'id',
      title: 'ID',
      dataIndex: 'id',
    },
    {
      key: 'email',
      title: '邮箱',
      dataIndex: 'email',
    },
    {
      key: 'domains',
      title: '域名',
      dataIndex: 'domains',
      render(col: string[]) {
        return (
          <Space>
            {col.map((item) => (
              <Tag key={item}>{item}</Tag>
            ))}
          </Space>
        );
      },
    },
    {
      key: 'status',
      title: '状态',
      dataIndex: 'status',
      render(col) {
        switch (col) {
          case 1:
            return <Tag color="green">已通过</Tag>;
          case 3:
            return <Tag color="red">已过期</Tag>;
          case 4:
            return <Tag color="blue">申请中</Tag>;
          case 5:
            return <Tag color="red">申请失败</Tag>;
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
            <Button
              type="text"
              style={{ color: 'var(--color-text-2)' }}
              iconOnly
              icon={<IconDownload />}
              onClick={() => {}}
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
      <Title heading={6}>证书管理</Title>
      <div className="flex flex-col">
        <div className="text-right">
          <CreateForm
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
              setVisible(true);
            }}
          >
            申请
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

export default Apply;

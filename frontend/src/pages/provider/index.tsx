import React from 'react';
import {
  Typography,
  Card,
  Table,
  Tag,
  Space,
  Button,
} from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table';
import { IconDelete, IconEdit, IconPlus } from '@arco-design/web-react/icon';
import Title from '@arco-design/web-react/es/Typography/title';

const columns: ColumnProps[] = [
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
  },
  {
    key: 'support',
    title: '支持',
    dataIndex: 'platform',
    render(col) {
      switch (col) {
        case 'tencent':
          return (
            <Space>
              <Tag color="green">CDN</Tag>
              <Tag color="arcoblue">DNS</Tag>
            </Space>
          );
      }
    },
  },
  {
    key: 'action',
    title: '操作',
    render(col) {
      return (
        <Space>
          <Button
            type="text"
            style={{ color: 'var(--color-text-2)' }}
            iconOnly
            icon={<IconEdit />}
          />
          <Button
            type="text"
            style={{ color: 'var(--color-text-2)' }}
            iconOnly
            icon={<IconDelete />}
          />
        </Space>
      );
    },
  },
];

const data = [
  {
    id: 1,
    name: '腾讯云',
    platform: 'tencent',
  },
];

function Provider() {
  return (
    <Card style={{ height: '80vh' }}>
      <Title heading={6}>厂商管理</Title>
      <div className="flex flex-col">
        <div className="text-right">
          <Button
            style={{ marginBottom: 10 }}
            type="primary"
            icon={<IconPlus />}
          >
            添加
          </Button>
        </div>
        <Table
          columns={columns}
          data={data}
          border={{
            wrapper: true,
          }}
        />
      </div>
    </Card>
  );
}

export default Provider;

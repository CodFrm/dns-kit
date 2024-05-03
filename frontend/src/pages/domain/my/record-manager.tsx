import {
  DomainItem,
  RecordExtraField,
  RecordItem,
  RecordTypes,
  useDomainDeleteMutation,
  useLazyRecordListQuery,
  useRecordCreateMutation,
  useRecordDeleteMutation,
  useRecordListQuery,
  useRecordUpdateMutation,
} from '@/services/domain.service';
import { useProviderEditMutation } from '@/services/provider.service';
import {
  Button,
  Form,
  FormInstance,
  Input,
  Message,
  Modal,
  Popconfirm,
  Select,
  Space,
  Switch,
  Table,
} from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';
import useForm from '@arco-design/web-react/es/Form/useForm';
import { ColumnProps } from '@arco-design/web-react/es/Table';
import { IconDelete, IconEdit, IconPlus } from '@arco-design/web-react/icon';
import React, { useContext } from 'react';
import { useEffect, useState } from 'react';

const EditableContext = React.createContext<{
  getForm?: () => FormInstance<RecordItem>;
}>({});

const RecordManager: React.FC<{
  visible: boolean;
  onOk: () => void;
  onCancel: () => void;
  domain?: DomainItem;
}> = (props) => {
  const [add, setAdd] = useState(false);
  const [lazyRecordList, { data, isLoading }] = useLazyRecordListQuery();
  const [deleteRecord, { data: deleteData, isLoading: deleteLoading }] =
    useRecordDeleteMutation();
  const [createRecord, { isLoading: createIsLoading }] =
    useRecordCreateMutation();
  const [updateRecord, { isLoading: updateIsLoading }] =
    useRecordUpdateMutation();
  const [recordList, setRecordList] = useState<RecordItem[]>([]);

  useEffect(() => {
    if (props.domain) {
      lazyRecordList(props.domain.id);
    }
  }, [props.domain]);

  const extraFields: RecordExtraField[] = data?.data?.extra || [];

  useEffect(() => {
    const recordList: RecordItem[] = data?.data?.list
      ? [...data.data.list]
      : [];
    if (add) {
      const extra = {};
      extraFields.forEach((field) => {
        extra[field.key] = field.default;
      });
      recordList.unshift({
        id: '',
        name: '',
        extra,
        type: 'A',
        value: '',
        ttl: 0,
        // @ts-ignore
        edit: true,
      });
    }
    setRecordList(recordList);
  }, [data, add]);

  const columns: ColumnProps<RecordItem>[] = [
    {
      key: 'name',
      title: '主机记录',
      dataIndex: 'name',
      editable: true,
      editRender(item) {
        return (
          <FormItem
            key="name"
            field="name"
            style={{ marginBottom: 0 }}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            initialValue={item.name}
            rules={[{ required: true }]}
          >
            <Input />
          </FormItem>
        );
      },
    },
    {
      key: 'type',
      title: '类型',
      dataIndex: 'type',
      editable: true,
      editRender(item) {
        return (
          <FormItem
            key="type"
            field="type"
            style={{ marginBottom: 0 }}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            initialValue={item.type}
            rules={[{ required: true }]}
          >
            <Select>
              {RecordTypes.map((type) => {
                return (
                  <Select.Option key={type} value={type}>
                    {type}
                  </Select.Option>
                );
              })}
            </Select>
          </FormItem>
        );
      },
    },
    {
      key: 'value',
      title: '记录值',
      dataIndex: 'value',
      editable: true,
      editRender(item) {
        return (
          <FormItem
            key="value"
            field="value"
            style={{ marginBottom: 0 }}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            initialValue={item.value}
            rules={[{ required: true }]}
          >
            <Input />
          </FormItem>
        );
      },
    },
    {
      key: 'ttl',
      title: 'TTL',
      dataIndex: 'ttl',
      editable: true,
      render(col) {
        return col === 1 ? '自动' : col;
      },
      editRender(item) {
        return (
          <FormItem
            key="ttl"
            field="ttl"
            style={{ marginBottom: 0 }}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            initialValue={item.ttl}
            rules={[{ required: true }]}
          >
            <Input type="number" />
          </FormItem>
        );
      },
    },
    ...extraFields.map((field2: any) => {
      const field = { ...field2 };
      field.dataIndex = 'extra.' + field.key;
      field.editable = true;
      const EditFormItem: React.FC<{
        initialValue: any;
        triggerPropName?: string;
        children: React.ReactNode;
      }> = ({ initialValue, triggerPropName, children }) => {
        return (
          <FormItem
            key={'extra.' + field.key}
            field={'extra.' + field.key}
            style={{ marginBottom: 0 }}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            initialValue={initialValue}
            rules={[{ required: true }]}
            triggerPropName={triggerPropName}
          >
            {children}
          </FormItem>
        );
      };
      switch (field.field_type) {
        case 'select':
          field.editRender = (item) => {
            return (
              <EditFormItem initialValue={item.extra[field.key]}>
                <Select>
                  {field.options.map((option) => (
                    <Select.Option key={option} value={option}>
                      {option}
                    </Select.Option>
                  ))}
                </Select>
              </EditFormItem>
            );
          };
          break;
        case 'switch':
          field.render = (col) => {
            return <Switch checked={col} />;
          };
          field.editRender = (item) => {
            return (
              <EditFormItem
                triggerPropName="checked"
                initialValue={item.extra[field.key]}
              >
                <Switch />
              </EditFormItem>
            );
          };
          break;
        default:
          field.editable = false;
          break;
      }
      return field;
    }),
    {
      key: 'action',
      title: '操作',
      render(_, item: RecordItem) {
        const { getForm } = useContext(EditableContext);
        return (
          <Space>
            <Button
              type="text"
              onClick={() => {
                setRecordList((list) => {
                  return list.map((record) => {
                    if (record.id === item.id) {
                      return {
                        ...record,
                        edit: true,
                      };
                    }
                    return record;
                  });
                });
              }}
            >
              编辑
            </Button>
            <Popconfirm
              title="确定删除吗？"
              onOk={() => {
                deleteRecord({
                  domain_id: props.domain.id,
                  record_id: item.id,
                })
                  .unwrap()
                  .then(() => {
                    Message.success('删除成功');
                  });
              }}
            >
              <Button type="text">删除</Button>
            </Popconfirm>
          </Space>
        );
      },
      editRender(item) {
        const { getForm } = useContext(EditableContext);
        return (
          <Space>
            <Button
              type="primary"
              loading={createIsLoading || updateIsLoading}
              onClick={() => {
                const values = getForm().getFieldsValue();
                setRecordList((list) => {
                  return list.map((record) => {
                    if (record.id === item.id) {
                      return {
                        ...record,
                        ...values,
                      };
                    }
                    return record;
                  });
                });
                if (item.id) {
                  // update
                  updateRecord({
                    domain_id: props.domain.id,
                    record_id: item.id,
                    record: {
                      type: values.type,
                      name: values.name,
                      value: values.value,
                      ttl: parseInt(values.ttl as unknown as string, 10),
                      extra: values.extra,
                    },
                  })
                    .unwrap()
                    .then(() => {
                      Message.success('修改成功');
                      setRecordList((list) => {
                        return list.map((record) => {
                          if (record.id === item.id) {
                            return {
                              ...record,
                              edit: false,
                            };
                          }
                          return record;
                        });
                      });
                    });
                } else {
                  // add
                  createRecord({
                    domain_id: props.domain.id,
                    record: {
                      type: values.type,
                      name: values.name,
                      value: values.value,
                      ttl: parseInt(values.ttl as unknown as string, 10),
                      extra: values.extra,
                    },
                  })
                    .unwrap()
                    .then(() => {
                      setAdd(false);
                      Message.success('添加成功');
                    });
                }
              }}
            >
              {item.id ? '修改' : '添加'}
            </Button>
            <Button
              type="text"
              onClick={() => {
                if (!item.id) {
                  setAdd(false);
                } else {
                  // 还原记录
                  setRecordList((list) => {
                    return list.map((record) => {
                      if (record.id === item.id) {
                        return data.data.list.find((r) => r.id === item.id);
                      }
                      return record;
                    });
                  });
                }
              }}
            >
              取消
            </Button>
          </Space>
        );
      },
    },
  ];

  const EditableRow = (props) => {
    const { children, record, className, ...rest } = props;
    const [form] = useForm<RecordItem>();
    const getForm = () => form;
    return (
      <EditableContext.Provider
        key={record.id}
        value={{
          getForm,
        }}
      >
        <Form
          style={{ display: 'table-row' }}
          children={children}
          form={form}
          wrapper="tr"
          wrapperProps={rest}
          className={`${className} editable-row`}
        ></Form>
      </EditableContext.Provider>
    );
  };

  const EditableCell = (props) => {
    const { children, className, rowData, column, onHandleSave } = props;

    return (
      <div
        key={column.key}
        className={'flex flex-row items-center gap-3 ' + className}
      >
        {rowData.edit && column.editRender
          ? column.editRender(rowData)
          : children}
      </div>
    );
  };

  return (
    <Modal
      title={'管理' + props.domain?.domain}
      style={{ width: '70%' }}
      visible={props.visible}
      confirmLoading={isLoading || deleteLoading}
      onOk={async () => {
        props.onOk();
      }}
      onCancel={() => props.onCancel()}
    >
      <div className="flex flex-col gap-3">
        <div className="text-right">
          <Button
            type="primary"
            icon={<IconPlus />}
            onClick={() => {
              setAdd(true);
            }}
          >
            添加
          </Button>
        </div>
        <Table
          columns={columns}
          components={{
            body: {
              row: EditableRow,
              cell: EditableCell,
            },
          }}
          data={recordList}
        />
      </div>
    </Modal>
  );
};

export default RecordManager;

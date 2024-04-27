import {
  DomainItem,
  RecordExtraField,
  RecordItem,
} from '@/services/domain.service';
import { useProviderEditMutation } from '@/services/provider.service';
import {
  Button,
  Form,
  FormInstance,
  Input,
  Modal,
  Select,
  Space,
  Table,
} from '@arco-design/web-react';
import FormItem from '@arco-design/web-react/es/Form/form-item';
import useForm from '@arco-design/web-react/es/Form/useForm';
import { ColumnProps } from '@arco-design/web-react/es/Table';
import { IconEdit, IconPlus } from '@arco-design/web-react/icon';
import React, { useContext } from 'react';
import { useEffect, useState } from 'react';

const tencent = [];

const EditableContext = React.createContext<{
  getForm?: () => FormInstance;
  edit?: boolean;
  setEdit?: (edit: boolean) => void;
}>({});

const extraFields: RecordExtraField[] = [
  {
    key: 'record_line',
    title: '线路',
    field_type: 'select',
    options: [
      '默认',
      '电信',
      '联通',
      '移动',
      '铁通',
      '广电',
      '教育网',
      '境内',
      '境外',
      '百度',
      '谷歌',
      '有道',
      '必应',
      '搜狗',
      '奇虎',
      '搜索引擎',
    ],
    default: '默认',
  },
];

const RecordManager: React.FC<{
  visible: boolean;
  onOk: () => void;
  onCancel: () => void;
  data?: DomainItem;
}> = (props) => {
  const [form] = Form.useForm();
  const [editProvider, { isLoading }] = useProviderEditMutation();
  const [add, setAdd] = useState(false);
  useEffect(() => {
    if (props.data) {
      form.setFieldsValue(props.data);
    } else {
      form.resetFields();
    }
  }, [props.data]);

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
              <Select.Option value="A">A</Select.Option>
              <Select.Option value="CNAME">CNAME</Select.Option>
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
    ...extraFields.map((field: any) => {
      field.dataIndex = 'extra.record_line';
      field.editable = true;
      const EditFormItem = ({ initialValue, children }) => {
        return (
          <FormItem
            key={field.key}
            field={field.key}
            style={{ marginBottom: 0 }}
            labelCol={{ span: 0 }}
            wrapperCol={{ span: 24 }}
            initialValue={initialValue}
            rules={[{ required: true }]}
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
        default:
          field.editable = false;
          break;
      }
      return field;
    }),
    {
      key: 'action',
      title: '操作',
      render(item) {
        return (
          <Space>
            <Button iconOnly type="text" icon={<IconPlus />} />
          </Space>
        );
      },
      editRender(item) {
        const { setEdit, getForm } = useContext(EditableContext);
        return (
          <Space>
            <Button
              type="primary"
              onClick={() => {
                console.log(getForm().getFieldsValue());
              }}
            >
              保存
            </Button>
            <Button
              type="text"
              onClick={() => {
                setEdit(false);
                if (!item.id) {
                  setAdd(false);
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
    const [edit, setEdit] = useState(record.id ? false : true);
    const [form] = useForm();
    const getForm = () => form;
    return (
      <EditableContext.Provider
        value={{
          getForm,
          edit: edit,
          setEdit: setEdit,
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
    const [show, setShow] = useState(false);
    const { getForm, edit, setEdit } = React.useContext(EditableContext);

    useEffect(() => {
      setShow(false);
    }, [edit]);

    return (
      <div
        onMouseMove={() => {
          setShow(true);
        }}
        onMouseLeave={() => {
          setShow(false);
        }}
        className={'flex flex-row items-center gap-3 ' + className}
      >
        {edit && column.editRender ? column.editRender(rowData) : children}
        {show && column.editable && !edit ? (
          <IconEdit
            className="cursor-pointer"
            onClick={() => {
              setEdit(true);
            }}
          />
        ) : (
          <div style={{ width: '32px', height: '32px' }} />
        )}
      </div>
    );
  };

  const data: RecordItem[] = [
    {
      id: '1',
      name: '2',
      value: '11',
      extra: {
        record_line: '默认',
      },
      type: 'A',
      ttl: 0,
    },
    {
      id: '1',
      name: '2',
      extra: {
        record_line: '默认',
      },
      type: 'A',
      value: '222',
      ttl: 0,
    },
  ];

  const recordList: RecordItem[] = data;
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
    });
  }

  return (
    <Modal
      title={'管理' + props.data?.name}
      style={{ width: '70%' }}
      visible={props.visible}
      confirmLoading={isLoading}
      onOk={async () => {
        form.validate().then((res) => {});
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

import React, { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import request from '../utils/request';

export default function VenuePage() {
  const [venues, setVenues] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form] = Form.useForm();

  const fetchVenues = async () => {
    setLoading(true);
    try {
      const res = await request.get('/api/venues');
      setVenues(res.data || []);
    } catch (e) {
      message.error('获取场地列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchVenues();
  }, []);

  const handleAdd = () => {
    setEditing(null);
    form.resetFields();
    setModalOpen(true);
  };

  const handleEdit = (record) => {
    setEditing(record);
    form.setFieldsValue(record);
    setModalOpen(true);
  };

  const handleDelete = async (id) => {
    try {
      await request.delete(`/api/venues/${id}`);
      message.success('删除成功');
      fetchVenues();
    } catch (e) {
      message.error(e.response?.data?.error || '删除失败');
    }
  };

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      if (editing) {
        await request.put(`/api/venues/${editing.id}`, values);
        message.success('修改成功');
      } else {
        await request.post('/api/venues', values);
        message.success('新增成功');
      }
      setModalOpen(false);
      fetchVenues();
    } catch (e) {
      if (e?.errorFields) {
        return;
      }
      message.error('保存失败');
    }
  };

  const columns = [
    { title: '场地名称', dataIndex: 'name' },
    { title: '地址', dataIndex: 'address' },
    { title: '负责人', dataIndex: 'manager_name' },
    { title: '联系方式', dataIndex: 'contact' },
    {
      title: '操作',
      dataIndex: 'action',
      render: (_, record) => (
        <>
          <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
          <Popconfirm title="确定删除该场地？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" danger>删除</Button>
          </Popconfirm>
        </>
      ),
    },
  ];

  return (
    <div>
      <Button type="primary" icon={<PlusOutlined />} style={{ marginBottom: 16 }} onClick={handleAdd}>
        新增场地
      </Button>
      <Table rowKey="id" columns={columns} dataSource={venues} loading={loading} />
      <Modal
        title={editing ? '编辑场地' : '新增场地'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleOk}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="场地名称" rules={[{ required: true, message: '请输入场地名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="address" label="地址">
            <Input />
          </Form.Item>
          <Form.Item name="manager_name" label="负责人">
            <Input />
          </Form.Item>
          <Form.Item name="contact" label="联系方式">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}


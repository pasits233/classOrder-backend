import React, { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, DatePicker, Select, message, Popconfirm } from 'antd';
import axios from 'axios';
import dayjs from 'dayjs';

export default function BookingPage() {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form] = Form.useForm();
  const [coaches, setCoaches] = useState([]);

  const fetchBookings = async () => {
    setLoading(true);
    try {
      const res = await axios.get('/api/bookings', {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
      });
      setBookings(res.data || []);
    } catch (e) {
      message.error('获取预约列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchCoaches = async () => {
    try {
      const res = await axios.get('/api/coaches', {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
      });
      setCoaches(res.data || []);
    } catch (e) {}
  };

  useEffect(() => {
    fetchBookings();
    fetchCoaches();
  }, []);

  const handleAdd = () => {
    setEditing(null);
    form.resetFields();
    setModalOpen(true);
  };

  const handleEdit = (record) => {
    setEditing(record);
    form.setFieldsValue({
      ...record,
      date: dayjs(record.date),
    });
    setModalOpen(true);
  };

  const handleDelete = async (id) => {
    try {
      await axios.delete(`/api/bookings/${id}`, {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
      });
      message.success('删除成功');
      fetchBookings();
    } catch (e) {
      message.error('删除失败');
    }
  };

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      let data = {
        ...values,
        date: values.date.format('YYYY-MM-DD'),
      };
      if (editing) {
        await axios.put(`/api/bookings/${editing.id}`, data, {
          headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
        });
        message.success('修改成功');
      } else {
        await axios.post('/api/bookings', data, {
          headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
        });
        message.success('添加成功');
      }
      setModalOpen(false);
      fetchBookings();
    } catch (e) {
      message.error('保存失败');
    }
  };

  const columns = [
    { title: '学员姓名', dataIndex: 'student_name' },
    { title: '教练', dataIndex: 'coach_id', render: id => coaches.find(c => c.id === id)?.name || '-' },
    { title: '日期', dataIndex: 'date' },
    { title: '时间段', dataIndex: 'time_slots' },
    { title: '操作', dataIndex: 'action', render: (_, record) => (
      <>
        <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
        <Popconfirm title="确定删除吗？" onConfirm={() => handleDelete(record.id)}>
          <Button type="link" danger>删除</Button>
        </Popconfirm>
      </>
    ) },
  ];

  return (
    <div>
      <Button type="primary" style={{ marginBottom: 16 }} onClick={handleAdd}>
        新增预约
      </Button>
      <Table rowKey="id" columns={columns} dataSource={bookings} loading={loading} />
      <Modal
        title={editing ? '编辑预约' : '新增预约'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleOk}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="student_name" label="学员姓名" rules={[{ required: true, message: '请输入学员姓名' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="coach_id" label="教练" rules={[{ required: true, message: '请选择教练' }]}>
            <Select>
              {coaches.map(coach => (
                <Select.Option key={coach.id} value={coach.id}>{coach.name}</Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="date" label="日期" rules={[{ required: true, message: '请选择日期' }]}>
            <DatePicker />
          </Form.Item>
          <Form.Item name="time_slots" label="时间段" rules={[{ required: true, message: '请输入时间段' }]}>
            <Input placeholder="如 09:00-09:30,10:00-10:30" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
} 
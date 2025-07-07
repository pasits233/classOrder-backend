import React, { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, Upload, message, Popconfirm, Image } from 'antd';
import { PlusOutlined, UploadOutlined } from '@ant-design/icons';
import axios from 'axios';

export default function CoachPage() {
  const [coaches, setCoaches] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form] = Form.useForm();
  const [uploading, setUploading] = useState(false);
  const [avatarUrl, setAvatarUrl] = useState('');

  const fetchCoaches = async () => {
    setLoading(true);
    try {
      const res = await axios.get('/api/coaches', {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
      });
      setCoaches(res.data || []);
    } catch (e) {
      message.error('获取教练列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCoaches();
  }, []);

  const handleAdd = () => {
    setEditing(null);
    setAvatarUrl('');
    form.resetFields();
    setModalOpen(true);
  };

  const handleEdit = (record) => {
    setEditing(record);
    setAvatarUrl(record.avatar_url || '');
    form.setFieldsValue({
      ...record,
      intro: record.description || '',
    });
    setModalOpen(true);
  };

  const handleDelete = async (id) => {
    try {
      await axios.delete(`/api/coaches/${id}`, {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
      });
      message.success('删除成功');
      fetchCoaches();
    } catch (e) {
      message.error('删除失败');
    }
  };

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      let data = {
        ...values,
        description: values.intro,
        avatar_url: avatarUrl,
      };
      delete data.intro;
      if (editing) {
        await axios.put(`/api/coaches/${editing.id}`, data, {
          headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
        });
        message.success('修改成功');
      } else {
        await axios.post('/api/coaches', data, {
          headers: { Authorization: 'Bearer ' + localStorage.getItem('token') },
        });
        message.success('添加成功');
      }
      setModalOpen(false);
      fetchCoaches();
    } catch (e) {
      message.error('保存失败');
    }
  };

  const handleUpload = async (info) => {
    setUploading(true);
    const formData = new FormData();
    formData.append('file', info.file);
    try {
      const res = await axios.post('/api/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          Authorization: 'Bearer ' + localStorage.getItem('token'),
        },
      });
      setAvatarUrl(res.data.url || res.data.file_url);
      message.success('上传成功');
    } catch (e) {
      message.error('上传失败');
    } finally {
      setUploading(false);
    }
  };

  const columns = [
    { title: '头像', dataIndex: 'avatar_url', render: url => url ? <Image src={url} width={40} /> : '-' },
    { title: '姓名', dataIndex: 'name' },
    { title: '简介', dataIndex: 'description' },
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
      <Button type="primary" icon={<PlusOutlined />} style={{ marginBottom: 16 }} onClick={handleAdd}>
        新增教练
      </Button>
      <Table rowKey="id" columns={columns} dataSource={coaches} loading={loading} />
      <Modal
        title={editing ? '编辑教练' : '新增教练'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleOk}
        confirmLoading={uploading}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="username" label="账号" rules={[{ required: true, message: '请输入账号' }]}> 
            <Input />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, message: '请输入密码' }]}> 
            <Input.Password />
          </Form.Item>
          <Form.Item name="name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]}> 
            <Input />
          </Form.Item>
          <Form.Item name="intro" label="简介" rules={[{ required: true, message: '请输入简介' }]}> 
            <Input />
          </Form.Item>
          <Form.Item label="头像">
            <Upload
              showUploadList={false}
              customRequest={({ file }) => handleUpload({ file })}
              accept="image/*"
            >
              <Button icon={<UploadOutlined />}>上传头像</Button>
            </Upload>
            {avatarUrl && <Image src={avatarUrl} width={60} style={{ marginTop: 8 }} />}
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
} 
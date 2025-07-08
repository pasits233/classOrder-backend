import React, { useEffect, useState } from 'react';
import { Card, Button, List, Modal, Form, Input, message, Spin } from 'antd';
import { Popconfirm } from 'antd';
import request from '../utils/request';
import { useNavigate } from 'react-router-dom';
import { logout, getRole } from '../utils/auth';
import './MobileCoachPage.css';

export default function MobileCoachPage() {
  const navigate = useNavigate();
  const role = getRole();
  const [coaches, setCoaches] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form] = Form.useForm();

  const fetchCoaches = async () => {
    setLoading(true);
    try {
      const res = await request.get('/api/coaches');
      setCoaches(res.data || []);
    } catch {
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
    form.resetFields();
    setModalOpen(true);
  };

  const handleEdit = (coach) => {
    setEditing(coach);
    form.setFieldsValue({
      ...coach,
      intro: coach.description || coach.intro || '',
    });
    setModalOpen(true);
  };

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      const data = { ...values, description: values.intro };
      delete data.intro;
      if (editing) {
        await request.put(`/api/coaches/${editing.id}`, data);
        message.success('修改成功');
      } else {
        await request.post('/api/coaches', data);
        message.success('添加成功');
      }
      setModalOpen(false);
      fetchCoaches();
    } catch {
      message.error('保存失败');
    }
  };

  return (
    <div className="mobile-coach-root">
      <div className="mobile-coach-navbar">
        <span className="mobile-coach-navbar-title">教练管理</span>
        <div className="mobile-coach-navbar-actions">
          <Button type="link" className="mobile-coach-navbar-link" onClick={() => navigate('/booking')}>预约管理</Button>
          <Button type="link" className="mobile-coach-navbar-link" onClick={() => { logout(); navigate('/login'); }}>退出登录</Button>
        </div>
      </div>
      <Button type="primary" block className="mobile-coach-add-btn" onClick={handleAdd}>
        新增教练
      </Button>
      {loading ? <Spin /> : (
        <List
          dataSource={coaches}
          renderItem={coach => (
            <Card className="mobile-coach-card" key={coach.id}>
              <div className="mobile-coach-card-actions">
                <Button type="link" className="mobile-coach-edit-btn" onClick={() => handleEdit(coach)}>
                  编辑
                </Button>
                <Popconfirm title="确定删除该教练吗？" onConfirm={async () => {
                  try {
                    await request.delete(`/api/coaches/${coach.id}`);
                    message.success('删除成功');
                    fetchCoaches();
                  } catch {
                    message.error('删除失败');
                  }
                }} okText="删除" cancelText="取消">
                  <Button type="link" danger className="mobile-coach-delete-btn">删除</Button>
                </Popconfirm>
              </div>
              <div className="mobile-coach-field"><span className="mobile-coach-label">姓名：</span>{coach.name}</div>
              <div className="mobile-coach-field"><span className="mobile-coach-label">简介：</span>{coach.description || coach.intro || '暂无简介'}</div>
              <div className="mobile-coach-field"><span className="mobile-coach-label">用户名：</span>{coach.username}</div>
            </Card>
          )}
        />
      )}
      <Modal
        title={editing ? '编辑教练' : '新增教练'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleOk}
        destroyOnClose
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]}> 
            <Input className="mobile-coach-input" />
          </Form.Item>
          <Form.Item name="intro" label="简介" rules={[{ required: true, message: '请输入简介' }]}> 
            <Input className="mobile-coach-input" />
          </Form.Item>
          <Form.Item name="username" label="用户名" rules={[{ required: true, message: '请输入用户名' }]}> 
            <Input className="mobile-coach-input" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: !editing, message: '请输入密码' }]}> 
            <Input.Password className="mobile-coach-input" placeholder={editing ? '如不修改可留空' : ''} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
} 
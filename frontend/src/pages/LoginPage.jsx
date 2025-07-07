import React, { useState } from 'react';
import { Form, Input, Button, Card, message, Select } from 'antd';
import { setRole } from '../utils/auth';
import { useNavigate } from 'react-router-dom';
import request from '../utils/request';

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values) => {
    setLoading(true);
    try {
      const res = await request.post('/api/login', {
        username: values.username,
        password: values.password,
        role: values.role,
      });
      if (res.data && res.data.token) {
        localStorage.setItem('token', res.data.token);
        setRole(values.role);
        message.success('登录成功');
        navigate('/');
      } else {
        message.error('登录失败');
      }
    } catch (e) {
      if (e.response && e.response.data && e.response.data.error) {
        message.error(e.response.data.error);
      } else {
        message.error('用户名或密码错误');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <Card title="后台登录" style={{ width: 350 }}>
        <Form onFinish={onFinish}>
          <Form.Item name="role" rules={[{ required: true, message: '请选择身份' }]}> 
            <Select placeholder="请选择身份">
              <Select.Option value="admin">管理员</Select.Option>
              <Select.Option value="coach">教练</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}> 
            <Input placeholder="用户名" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}> 
            <Input.Password placeholder="密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
} 
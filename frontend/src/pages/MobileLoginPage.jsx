import React, { useState } from 'react';
import { Form, Input, Button, message, Card, Select } from 'antd';
import request from '../utils/request';
import { setRole } from '../utils/auth';
import { useNavigate } from 'react-router-dom';
import './MobileLoginPage.css';

export default function MobileLoginPage() {
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
        navigate(values.role === 'admin' ? '/coach' : '/booking');
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
    <div className="mobile-login-root">
      <Card className="mobile-login-card">
        <div className="mobile-login-title">课程预约系统</div>
        <Form layout="vertical" onFinish={onFinish} autoComplete="off">
          <Form.Item name="role" label="角色" rules={[{ required: true, message: '请选择角色' }]}> 
            <Select className="mobile-login-input" placeholder="请选择身份">
              <Select.Option value="admin">管理员</Select.Option>
              <Select.Option value="coach">教练</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="username" label="用户名" rules={[{ required: true, message: '请输入用户名' }]}> 
            <Input className="mobile-login-input" placeholder="请输入用户名" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, message: '请输入密码' }]}> 
            <Input.Password className="mobile-login-input" placeholder="请输入密码" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading} className="mobile-login-btn">
            登录
          </Button>
        </Form>
      </Card>
    </div>
  );
} 
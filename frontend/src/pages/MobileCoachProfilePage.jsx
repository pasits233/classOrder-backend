import React, { useEffect, useState } from 'react';
import { Form, Input, Button, message, Upload, Avatar, Spin } from 'antd';
import { UploadOutlined, UserOutlined, LogoutOutlined } from '@ant-design/icons';
import request from '../utils/request';
import { logout } from '../utils/auth';

export default function MobileCoachProfilePage() {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [profile, setProfile] = useState(null);
  const [avatarUrl, setAvatarUrl] = useState('');
  const [uploading, setUploading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmNewPassword, setConfirmNewPassword] = useState('');

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    setLoading(true);
    try {
      const res = await request.get('/api/coach/profile');
      setProfile(res.data);
      setAvatarUrl(res.data.avatar_url || '');
      form.setFieldsValue({
        name: res.data.name,
        description: res.data.description,
        username: res.data.username,
        password: '',
      });
    } catch {
      message.error('获取个人信息失败');
    } finally {
      setLoading(false);
    }
  };

  const handleUpload = async ({ file }) => {
    setUploading(true);
    const formData = new FormData();
    formData.append('file', file);
    try {
      const res = await request.post('/api/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      setAvatarUrl(res.data.file_url || res.data.url);
      message.success('头像上传成功');
    } catch {
      message.error('头像上传失败');
    } finally {
      setUploading(false);
    }
  };

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      if (newPassword || confirmNewPassword) {
        if (!oldPassword) {
          message.error('请输入原密码');
          return;
        }
        if (newPassword !== confirmNewPassword) {
          message.error('两次输入的新密码不一致');
          return;
        }
      }
      setSaving(true);
      await request.put('/api/coach/profile', {
        name: values.name,
        description: values.description,
        avatar_url: avatarUrl,
        old_password: oldPassword || undefined,
        password: newPassword || undefined,
      });
      message.success('保存成功');
      fetchProfile();
      setOldPassword('');
      setNewPassword('');
      setConfirmNewPassword('');
      form.setFieldValue('password', '');
    } catch {
      message.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  if (loading || !profile) return <Spin style={{ marginTop: 80 }} />;

  return (
    <div style={{ maxWidth: 400, margin: '0 auto', padding: 16 }}>
      <div style={{ textAlign: 'center', marginBottom: 24 }}>
        <Avatar size={80} src={avatarUrl} icon={<UserOutlined />} />
        <div style={{ marginTop: 8 }}>
          <Upload showUploadList={false} customRequest={handleUpload} accept="image/*">
            <Button icon={<UploadOutlined />} loading={uploading} size="small">更换头像</Button>
          </Upload>
        </div>
      </div>
      <Form form={form} layout="vertical" autoComplete="off">
        <Form.Item label="用户名" name="username">
          <Input disabled />
        </Form.Item>
        <Form.Item label="姓名" name="name" rules={[{ required: true, message: '请输入姓名' }]}>
          <Input placeholder="请输入姓名" />
        </Form.Item>
        <Form.Item label="简介" name="description">
          <Input.TextArea rows={3} placeholder="请输入简介" />
        </Form.Item>
        <Form.Item label="原密码">
          <Input.Password value={oldPassword} onChange={e => setOldPassword(e.target.value)} autoComplete="current-password" />
        </Form.Item>
        <Form.Item label="新密码">
          <Input.Password value={newPassword} onChange={e => setNewPassword(e.target.value)} autoComplete="new-password" />
        </Form.Item>
        <Form.Item label="确认新密码">
          <Input.Password value={confirmNewPassword} onChange={e => setConfirmNewPassword(e.target.value)} autoComplete="new-password" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" block onClick={handleSave} loading={saving}>保存</Button>
        </Form.Item>
        <Form.Item>
          <Button icon={<LogoutOutlined />} block danger onClick={() => { logout(); window.location.href = '/login'; }}>退出登录</Button>
        </Form.Item>
      </Form>
    </div>
  );
} 
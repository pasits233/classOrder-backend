import React, { useEffect, useState } from 'react';
import { Table, Button, Modal, Form, Input, InputNumber, Upload, message, Popconfirm, Image } from 'antd';
import { PlusOutlined, UploadOutlined } from '@ant-design/icons';
import request from '../utils/request';

export default function CoachPage() {
  const [coaches, setCoaches] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form] = Form.useForm();
  const [uploading, setUploading] = useState(false);
  const [avatarUrl, setAvatarUrl] = useState('');
  const [videoModalOpen, setVideoModalOpen] = useState(false);
  const [videoEditorOpen, setVideoEditorOpen] = useState(false);
  const [currentCoach, setCurrentCoach] = useState(null);
  const [coachVideos, setCoachVideos] = useState([]);
  const [videoLoading, setVideoLoading] = useState(false);
  const [editingVideo, setEditingVideo] = useState(null);
  const [videoForm] = Form.useForm();
  const [videoUploading, setVideoUploading] = useState(false);
  const [coverUploading, setCoverUploading] = useState(false);

  const fetchCoaches = async () => {
    setLoading(true);
    try {
      const res = await request.get('/api/coaches');
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
      await request.delete(`/api/coaches/${id}`);
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
        await request.put(`/api/coaches/${editing.id}`, data);
        message.success('修改成功');
      } else {
        await request.post('/api/coaches', data);
        message.success('添加成功，默认密码为: coach123');
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
      const res = await request.post('/api/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
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

  const fetchCoachVideos = async (coachId) => {
    setVideoLoading(true);
    try {
      const res = await request.get(`/api/coaches/${coachId}/videos`);
      setCoachVideos(res.data || []);
    } catch (e) {
      message.error('获取教练视频失败');
    } finally {
      setVideoLoading(false);
    }
  };

  const openVideoModal = (coach) => {
    setCurrentCoach(coach);
    setVideoModalOpen(true);
    fetchCoachVideos(coach.id);
  };

  const handleAddVideo = () => {
    setEditingVideo(null);
    videoForm.resetFields();
    setVideoEditorOpen(true);
  };

  const handleEditVideo = (video) => {
    setEditingVideo(video);
    videoForm.setFieldsValue(video);
    setVideoEditorOpen(true);
  };

  const handleSaveVideo = async () => {
    if (!currentCoach) return;
    try {
      const values = await videoForm.validateFields();
      const payload = {
        ...values,
        sort_order: values.sort_order ?? 0,
      };
      if (editingVideo) {
        await request.put(`/api/coach-videos/${editingVideo.id}`, payload);
        message.success('视频更新成功');
      } else {
        await request.post(`/api/coaches/${currentCoach.id}/videos`, payload);
        message.success('视频添加成功');
      }
      setVideoEditorOpen(false);
      fetchCoachVideos(currentCoach.id);
    } catch (e) {
      if (e?.errorFields) return;
      message.error('保存视频失败');
    }
  };

  const handleDeleteVideo = async (videoId) => {
    try {
      await request.delete(`/api/coach-videos/${videoId}`);
      message.success('删除成功');
      if (currentCoach) {
        fetchCoachVideos(currentCoach.id);
      }
    } catch (e) {
      message.error('删除失败');
    }
  };

  const handleVideoUpload = async ({ file, onSuccess, onError }) => {
    setVideoUploading(true);
    const formData = new FormData();
    formData.append('file', file);
    try {
      const res = await request.post('/api/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      const url = res.data.url || res.data.file_url;
      videoForm.setFieldsValue({ video_url: url });
      if (onSuccess) onSuccess(res.data);
      message.success('视频上传成功');
    } catch (e) {
      if (onError) onError(e);
      message.error('视频上传失败');
    } finally {
      setVideoUploading(false);
    }
  };

  const handleCoverUpload = async ({ file, onSuccess, onError }) => {
    setCoverUploading(true);
    const formData = new FormData();
    formData.append('file', file);
    try {
      const res = await request.post('/api/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      const url = res.data.url || res.data.file_url;
      videoForm.setFieldsValue({ cover_url: url });
      if (onSuccess) onSuccess(res.data);
      message.success('封面上传成功');
    } catch (e) {
      if (onError) onError(e);
      message.error('封面上传失败');
    } finally {
      setCoverUploading(false);
    }
  };

  const columns = [
    { title: '头像', dataIndex: 'avatar_url', render: url => url ? <Image src={url} width={40} /> : '-' },
    { title: '姓名', dataIndex: 'name' },
    { title: '简介', dataIndex: 'description' },
    { title: '操作', dataIndex: 'action', render: (_, record) => (
      <>
        <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
        <Button type="link" onClick={() => handleResetPassword(record.id)}>重置密码</Button>
        <Button type="link" onClick={() => openVideoModal(record)}>管理视频</Button>
        <Popconfirm title="确定删除吗？" onConfirm={() => handleDelete(record.id)}>
          <Button type="link" danger>删除</Button>
        </Popconfirm>
      </>
    ) },
  ];

  const videoColumns = [
    { title: '标题', dataIndex: 'title' },
    { title: '排序', dataIndex: 'sort_order' },
    { title: '封面', dataIndex: 'cover_url', render: url => url ? <Image src={url} width={80} /> : '-' },
    { title: '视频链接', dataIndex: 'video_url', render: url => url ? <a href={url} target="_blank" rel="noreferrer">查看</a> : '-' },
    {
      title: '操作',
      dataIndex: 'action',
      render: (_, record) => (
        <>
          <Button type="link" onClick={() => handleEditVideo(record)}>编辑</Button>
          <Popconfirm title="确定删除该视频？" onConfirm={() => handleDeleteVideo(record.id)}>
            <Button type="link" danger>删除</Button>
          </Popconfirm>
        </>
      ),
    },
  ];

  const handleResetPassword = async (coachId) => {
    try {
      await request.post(`/api/coaches/${coachId}/reset-password`);
      message.success('密码重置成功，新密码为: coach123');
    } catch (e) {
      message.error('密码重置失败');
    }
  };

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
      <Modal
        title={currentCoach ? `${currentCoach.name} - 视频库` : '教练视频'}
        open={videoModalOpen}
        onCancel={() => {
          setVideoModalOpen(false);
          setVideoEditorOpen(false);
        }}
        footer={null}
        width={840}
      >
        <Button type="primary" style={{ marginBottom: 16 }} onClick={handleAddVideo}>
          新增视频
        </Button>
        <Table
          rowKey="id"
          columns={videoColumns}
          dataSource={coachVideos}
          loading={videoLoading}
          pagination={false}
        />
      </Modal>
      <Modal
        title={editingVideo ? '编辑视频' : '新增视频'}
        open={videoEditorOpen}
        onCancel={() => setVideoEditorOpen(false)}
        onOk={handleSaveVideo}
        okText="保存"
        cancelText="取消"
        destroyOnClose
      >
        <Form form={videoForm} layout="vertical">
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="video_url" label="视频地址" rules={[{ required: true, message: '请上传或填写视频地址' }]}>
            <Input addonAfter={
              <Upload
                showUploadList={false}
                customRequest={handleVideoUpload}
                accept="video/*"
              >
                <Button loading={videoUploading}>上传</Button>
              </Upload>
            } />
          </Form.Item>
          <Form.Item name="cover_url" label="封面地址">
            <Input addonAfter={
              <Upload
                showUploadList={false}
                customRequest={handleCoverUpload}
                accept="image/*"
              >
                <Button loading={coverUploading}>上传</Button>
              </Upload>
            } />
          </Form.Item>
          <Form.Item name="sort_order" label="排序">
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
} 
import React, { useEffect, useState } from 'react';
import { Card, Button, List, Modal, Form, Input, message, Spin } from 'antd';
import request from '../utils/request';
import './MobileVenuePage.css';

export default function MobileVenuePage() {
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
      message.error('获取场地失败');
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
      if (e?.errorFields) return;
      message.error('保存失败');
    }
  };

  const handleDelete = (record) => {
    Modal.confirm({
      title: '确定删除该场地？',
      onConfirm: async () => {
        try {
          await request.delete(`/api/venues/${record.id}`);
          message.success('删除成功');
          fetchVenues();
        } catch (err) {
          message.error(err.response?.data?.error || '删除失败');
        }
      },
    });
  };

  return (
    <div className="mobile-venue-root">
      <Button type="primary" block className="mobile-venue-add-btn" onClick={handleAdd}>
        新增场地
      </Button>
      {loading ? <Spin /> : (
        <List
          dataSource={venues}
          renderItem={venue => (
            <Card className="mobile-venue-card" key={venue.id}>
              <div className="mobile-venue-field"><span className="mobile-venue-label">名称：</span>{venue.name}</div>
              <div className="mobile-venue-field"><span className="mobile-venue-label">地址：</span>{venue.address || '—'}</div>
              <div className="mobile-venue-field"><span className="mobile-venue-label">负责人：</span>{venue.manager_name || '—'}</div>
              <div className="mobile-venue-field"><span className="mobile-venue-label">联系方式：</span>{venue.contact || '—'}</div>
              <div className="mobile-venue-actions">
                <Button type="link" onClick={() => handleEdit(venue)}>编辑</Button>
                <Button type="link" danger onClick={() => handleDelete(venue)}>删除</Button>
              </div>
            </Card>
          )}
        />
      )}
      <Modal
        title={editing ? '编辑场地' : '新增场地'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleOk}
        okText="保存"
        cancelText="取消"
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


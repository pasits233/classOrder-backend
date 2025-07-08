import React, { useEffect, useState } from 'react';
import { Button, Modal, Form, Input, DatePicker, Select, message, Tag, List, Drawer, Spin } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import request from '../utils/request';
import { getRole, getUserId } from '../utils/auth';
import './MobileBookingPage.css';
import { useNavigate } from 'react-router-dom';

const TIME_SLOTS = [];
for (let hour = 9; hour < 18; hour++) {
  TIME_SLOTS.push(`${hour.toString().padStart(2, '0')}:00-${hour.toString().padStart(2, '0')}:30`);
  TIME_SLOTS.push(`${hour.toString().padStart(2, '0')}:30-${(hour + 1).toString().padStart(2, '0')}:00`);
}

export default function MobileBookingPage() {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [coaches, setCoaches] = useState([]);
  const [selectedCoach, setSelectedCoach] = useState(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form] = Form.useForm();
  const [selectedSlots, setSelectedSlots] = useState([]);
  const [unavailableSlots, setUnavailableSlots] = useState([]);
  const [selectedDate, setSelectedDate] = useState(null);
  const role = getRole();
  const userId = getUserId();
  const navigate = useNavigate();

  // 获取教练列表
  const fetchCoaches = async () => {
    try {
      const res = await request.get('/api/coaches');
      setCoaches(res.data || []);
      // 管理员不再默认选中第一个教练，selectedCoach 保持 null
      if (role === 'coach') {
        const myCoach = res.data.find(c => String(c.user_id) === String(userId));
        if (myCoach) setSelectedCoach(myCoach.id);
      }
    } catch (e) {
      message.error('获取教练列表失败');
    }
  };

  // 获取预约列表
  const fetchBookings = async (coachId, date) => {
    setLoading(true);
    try {
      const params = {};
      if (role === 'admin') {
        if (coachId !== null && coachId !== undefined) params.coach_id = coachId;
        if (date) params.date = dayjs(date).format('YYYY-MM-DD');
      } else {
        params.coach_id = coachId;
        if (date) params.date = dayjs(date).format('YYYY-MM-DD');
      }
      const res = await request.get('/api/bookings', { params });
      setBookings(Array.isArray(res.data) ? res.data : []);
    } catch (e) {
      message.error('获取预约列表失败');
      setBookings([]);
    } finally {
      setLoading(false);
    }
  };

  // 查询某教练某天已预约的所有时间段
  const fetchUnavailableSlots = async (coachId, date, editingId) => {
    if (!coachId || !date) {
      setUnavailableSlots([]);
      return;
    }
    try {
      const res = await request.get('/api/bookings', {
        params: { coach_id: coachId, date: typeof date === 'string' ? date : date.format('YYYY-MM-DD') },
      });
      let slots = [];
      res.data.forEach(b => {
        if (!editingId || b.id !== editingId) {
          slots = slots.concat(b.time_slots.split(',').map(s => s.trim()));
        }
      });
      setUnavailableSlots(slots);
    } catch {
      setUnavailableSlots([]);
    }
  };

  useEffect(() => {
    fetchCoaches();
  }, []);

  // 管理员进入时默认显示全部预约，筛选器变动自动刷新
  useEffect(() => {
    if (role === 'admin') {
      fetchBookings(selectedCoach, selectedDate);
    } else if (selectedCoach) {
      fetchBookings(selectedCoach, selectedDate);
    }
  }, [role, selectedCoach, selectedDate]);

  // useEffect: 管理员首次进入selectedCoach为null
  useEffect(() => {
    if (role === 'admin' && selectedCoach === null) {
      fetchBookings(null, selectedDate);
    }
  }, [role, selectedCoach, selectedDate]);

  // 新增/编辑预约
  const handleAdd = () => {
    setEditing(null);
    form.resetFields();
    setSelectedSlots([]);
    setDrawerOpen(true);
    fetchUnavailableSlots(selectedCoach, selectedDate);
  };

  const handleEdit = (record) => {
    setEditing(record);
    form.setFieldsValue({
      ...record,
      date: dayjs(record.date),
    });
    setSelectedSlots(record.time_slots ? record.time_slots.split(',').map(s => s.trim()) : []);
    setDrawerOpen(true);
    fetchUnavailableSlots(record.coach_id, dayjs(record.date), record.id);
  };

  // 时间段按钮点击
  const handleSlotClick = (slot) => {
    setSelectedSlots(prev =>
      prev.includes(slot)
        ? prev.filter(s => s !== slot)
        : [...prev, slot]
    );
  };

  // 提交
  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      let data = {
        ...values,
        date: values.date.format('YYYY-MM-DD'),
        time_slots: selectedSlots.join(', '),
      };
      if (role === 'coach') {
        const myCoach = coaches.find(c => String(c.user_id) === String(userId));
        if (myCoach) data.coach_id = myCoach.id;
      }
      if (editing) {
        await request.put(`/api/bookings/${editing.id}`, data);
        message.success('修改成功');
      } else {
        await request.post('/api/bookings', data);
        message.success('添加成功');
      }
      setDrawerOpen(false);
      fetchBookings(selectedCoach, selectedDate);
    } catch (e) {
      message.error('保存失败');
    }
  };

  // 日期变更
  const handleDateChange = (date) => {
    setSelectedDate(date);
    fetchUnavailableSlots(selectedCoach, date, editing?.id);
    setSelectedSlots([]);
  };

  // 教练变更
  const handleCoachChange = (coachId) => {
    setSelectedCoach(coachId);
    fetchUnavailableSlots(coachId, selectedDate, editing?.id);
    setSelectedSlots([]);
  };

  // 在顶部加个人中心按钮（仅教练）
  const profileBtn = role === 'coach' ? (
    <Button type="link" style={{ float: 'right', marginTop: 8 }} onClick={() => navigate('/profile')}>个人中心</Button>
  ) : null;

  // 顶部筛选器
  const filterBar = (
    <div className="mobile-booking-header">
      {profileBtn}
      <DatePicker
        value={selectedDate}
        onChange={handleDateChange}
        className="mobile-booking-date"
        allowClear
        style={{ width: role === 'admin' ? '48%' : '100%' }}
        placeholder="选择日期"
      />
      {role === 'admin' && (
        <Select
          value={selectedCoach}
          onChange={value => setSelectedCoach(value)}
          className="mobile-booking-coach"
          style={{ width: '48%' }}
          allowClear
          placeholder="全部教练"
        >
          <Select.Option value={null}>全部教练</Select.Option>
          {coaches.map(coach => (
            <Select.Option key={coach.id} value={coach.id}>{coach.name}</Select.Option>
          ))}
        </Select>
      )}
    </div>
  );

  // 在组件内补充 handleDelete 方法：
  const handleDelete = async (record) => {
    Modal.confirm({
      title: '确认删除该预约？',
      onOk: async () => {
        try {
          await request.delete(`/api/bookings/${record.id}`);
          message.success('删除成功');
          fetchBookings(selectedCoach, selectedDate);
        } catch {
          message.error('删除失败');
        }
      },
    });
  };

  return (
    <div className="mobile-booking-root">
      {filterBar}
      <div className="mobile-booking-list">
        {loading ? <Spin /> : (
          bookings.length === 0 ? (
            <div className="mobile-booking-empty">暂无预约</div>
          ) : (
            <List
              dataSource={bookings}
              renderItem={item => (
                <div className="mobile-booking-card">
                  <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8, width: '100%', gap: 0 }}>
                    <div style={{ flex: 1, textAlign: 'center' }}>
                      <Button
                        type="text"
                        block
                        style={{ color: '#1677ff', fontWeight: 600, fontSize: 16 }}
                        onClick={() => handleEdit(item)}
                      >
                        编辑
                      </Button>
                    </div>
                    <div style={{ flex: 1, textAlign: 'center' }}>
                      <Button
                        type="text"
                        danger
                        block
                        style={{ fontWeight: 600, fontSize: 16 }}
                        onClick={() => handleDelete(item)}
                      >
                        删除
                      </Button>
                    </div>
                  </div>
                  <div className="mobile-booking-field"><span className="mobile-booking-label">日期：</span>{item.date}</div>
                  <div className="mobile-booking-field"><span className="mobile-booking-label">学员：</span>{item.student_name}</div>
                  <div className="mobile-booking-field"><span className="mobile-booking-label">时间段：</span><Tag color="blue">{item.time_slots}</Tag></div>
                  <div className="mobile-booking-field"><span className="mobile-booking-label">教练：</span>{coaches.find(c => c.id === item.coach_id)?.name || '-'}</div>
                </div>
              )}
            />
          )
        )}
      </div>
      <Button
        type="primary"
        shape="circle"
        icon={<PlusOutlined />}
        className="mobile-booking-add-btn"
        onClick={handleAdd}
        size="large"
      />
      <Drawer
        title={editing ? '编辑预约' : '新增预约'}
        placement="bottom"
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        height="80vh"
        className="mobile-booking-drawer"
        bodyStyle={{ padding: 16 }}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="student_name" label="学员姓名" rules={[{ required: true, message: '请输入学员姓名' }]}> 
            <Input />
          </Form.Item>
          {role === 'admin' && (
            <Form.Item name="coach_id" label="教练" rules={[{ required: true, message: '请选择教练' }]}> 
              <Select>
                {coaches.map(coach => (
                  <Select.Option key={coach.id} value={coach.id}>{coach.name}</Select.Option>
                ))}
              </Select>
            </Form.Item>
          )}
          <Form.Item name="date" label="日期" rules={[{ required: true, message: '请选择日期' }]}> 
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item label="时间段" required>
            <div className="mobile-booking-slots">
              {TIME_SLOTS.map(slot => (
                <Button
                  key={slot}
                  type={selectedSlots.includes(slot) ? 'primary' : 'default'}
                  onClick={() => unavailableSlots.includes(slot) ? null : handleSlotClick(slot)}
                  disabled={unavailableSlots.includes(slot)}
                  className="mobile-booking-slot-btn"
                >
                  {slot}
                </Button>
              ))}
            </div>
            {selectedSlots.length === 0 && <div style={{ color: 'red', marginTop: 4 }}>请选择至少一个时间段</div>}
          </Form.Item>
        </Form>
        <Button type="primary" block onClick={handleOk} style={{ marginTop: 16 }}>
          保存
        </Button>
      </Drawer>
    </div>
  );
} 
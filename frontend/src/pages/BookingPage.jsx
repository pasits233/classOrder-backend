import React, { useEffect, useState } from 'react';
import { Card, Button, Select, message, Spin, List, Tag, Modal, Form, Input, DatePicker } from 'antd';
import request from '../utils/request';
import dayjs from 'dayjs';
import { getRole, getUserId } from '../utils/auth';
import './BookingPage.css';

const TIME_SLOTS = [];
for (let hour = 9; hour < 18; hour++) {
  const start = `${hour.toString().padStart(2, '0')}:00`;
  const end = `${(hour + 1).toString().padStart(2, '0')}:00`;
  TIME_SLOTS.push(`${start}-${end}`);
}

export default function BookingPage() {
  // useState 声明全部提前
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [coaches, setCoaches] = useState([]);
  const [venues, setVenues] = useState([]);
  const [selectedCoach, setSelectedCoach] = useState(null);
  const [selectedVenue, setSelectedVenue] = useState(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form] = Form.useForm();
  const [selectedSlots, setSelectedSlots] = useState([]);
  const [unavailableSlots, setUnavailableSlots] = useState([]);
  const [selectedDate, setSelectedDate] = useState(null);
  const role = getRole();
  const userId = getUserId();

  // 获取教练列表
  const fetchCoaches = async () => {
    try {
      const res = await request.get('/api/coaches');
      setCoaches(res.data || []);
      // 管理员默认选第一个教练
      if (role === 'admin' && res.data && res.data.length > 0) {
        setSelectedCoach(res.data[0].id);
      }
      // 教练默认选自己
      if (role === 'coach') {
        // 通过 userId 找到自己的 coach.id，类型统一为字符串比较
        const myCoach = res.data.find(c => String(c.user_id) === String(userId));
        if (myCoach) setSelectedCoach(myCoach.id);
      }
    } catch (e) {
      message.error('获取教练列表失败');
    }
  };

  const fetchVenues = async () => {
    try {
      const res = await request.get('/api/venues');
      setVenues(res.data || []);
    } catch (e) {
      message.error('获取场地列表失败');
    }
  };

  // 获取预约列表
  const fetchBookings = async (coachId, date, venueId) => {
    setLoading(true);
    try {
      const params = {};
      if (role === 'admin') {
        if (coachId) params.coach_id = coachId;
        if (venueId) params.venue_id = venueId;
        if (date) params.date = date.format('YYYY-MM-DD');
      } else {
        params.coach_id = coachId;
        if (venueId) params.venue_id = venueId;
        if (date) params.date = date.format('YYYY-MM-DD');
      }
      const res = await request.get('/api/bookings', { params });
      setBookings(res.data || []);
    } catch (e) {
      message.error('获取预约列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCoaches();
    fetchVenues();
  }, []);

  useEffect(() => {
    if (selectedCoach) {
      fetchBookings(selectedCoach, selectedDate, selectedVenue);
    } else {
      fetchBookings(null, selectedDate, selectedVenue);
    }
  }, [selectedCoach, selectedDate, selectedVenue]);

  // 以日期分组
  const bookingsByDate = bookings.reduce((acc, cur) => {
    const date = cur.date;
    if (!acc[date]) acc[date] = [];
    acc[date].push(cur);
    return acc;
  }, {});

  // 顶部筛选器
  const filterBar = (
    <div style={{ display: 'flex', gap: 16, marginBottom: 16 }}>
      <DatePicker
        value={selectedDate}
        onChange={setSelectedDate}
        style={{ width: 180 }}
        placeholder="选择日期"
        allowClear
      />
      {role === 'admin' && (
        <Select
          style={{ width: 200 }}
          value={selectedCoach}
          onChange={setSelectedCoach}
          allowClear
          placeholder="全部教练"
        >
          <Select.Option value={null}>全部教练</Select.Option>
          {coaches.map(coach => (
            <Select.Option key={coach.id} value={coach.id}>{coach.name}</Select.Option>
          ))}
        </Select>
      )}
      <Select
        style={{ width: 200 }}
        value={selectedVenue}
        onChange={value => setSelectedVenue(value || null)}
        allowClear
        placeholder="全部场地"
      >
        <Select.Option value={null}>全部场地</Select.Option>
        {venues.map(venue => (
          <Select.Option key={venue.id} value={venue.id}>{venue.name}</Select.Option>
        ))}
      </Select>
    </div>
  );

  // 查询某教练某天已预约的所有时间段
  const fetchUnavailableSlots = async (coachId, date, venueId, editingId) => {
    if (!coachId || !date || !venueId) {
      setUnavailableSlots([]);
      return;
    }
    try {
      const res = await request.get('/api/bookings', {
        params: { 
          coach_id: coachId, 
          venue_id: venueId,
          date: typeof date === 'string' ? date : date.format('YYYY-MM-DD') 
        },
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

  const resolveCoachIdForForm = () => {
    if (role === 'coach') {
      const myCoach = coaches.find(c => String(c.user_id) === String(userId));
      return myCoach ? myCoach.id : null;
    }
    return form.getFieldValue('coach_id') || selectedCoach;
  };

  const refreshUnavailableSlots = (overrides = {}) => {
    const coachId = overrides.coachId ?? resolveCoachIdForForm();
    const dateVal = overrides.date ?? form.getFieldValue('date');
    const venueId = overrides.venueId ?? form.getFieldValue('venue_id');
    const editingId = overrides.editingId ?? (editing ? editing.id : undefined);

    if (coachId && dateVal && venueId) {
      const dateStr = typeof dateVal === 'string' ? dateVal : dateVal.format('YYYY-MM-DD');
      fetchUnavailableSlots(coachId, dateStr, venueId, editingId);
    } else {
      setUnavailableSlots([]);
    }
  };

  // 新增/编辑预约
  const handleAdd = () => {
    setEditing(null);
    form.resetFields();
    setSelectedSlots([]);
    setUnavailableSlots([]);
    setModalOpen(true);
  };

  const handleEdit = (record) => {
    setEditing(record);
    form.setFieldsValue({
      ...record,
      date: dayjs(record.date),
      venue_id: record.venue_id,
    });
    setSelectedSlots(record.time_slots ? record.time_slots.split(',').map(s => s.trim()) : []);
    setModalOpen(true);
    refreshUnavailableSlots({
      coachId: record.coach_id,
      date: dayjs(record.date),
      venueId: record.venue_id,
      editingId: record.id,
    });
  };

  const handleCoachChangeInModal = (value) => {
    form.setFieldsValue({ coach_id: value });
    refreshUnavailableSlots({ coachId: value });
    setSelectedSlots([]);
  };

  const handleDateChangeInModal = (value) => {
    refreshUnavailableSlots({ date: value });
    setSelectedSlots([]);
  };

  const handleVenueChangeInModal = (value) => {
    form.setFieldsValue({ venue_id: value });
    refreshUnavailableSlots({ venueId: value });
    setSelectedSlots([]);
  };

  // 时间段按钮点击
  const handleSlotClick = (slot) => {
    setSelectedSlots(prev =>
      prev.includes(slot)
        ? prev.filter(s => s !== slot)
        : [...prev, slot]
    );
  };

  // 提交时自动填充 time_slots 字段
  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      let data = {
        ...values,
        date: values.date.format('YYYY-MM-DD'),
        time_slots: selectedSlots.join(', '),
      };
      if (role === 'coach') {
        // 通过 userId 找到自己的 coach.id
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
      setModalOpen(false);
      fetchBookings(selectedCoach, selectedDate, selectedVenue);
    } catch (e) {
      message.error('保存失败');
    }
  };

  return (
    <div>
      {filterBar}
      <Button type="primary" style={{ marginBottom: 16 }} onClick={handleAdd}>
        新增预约
      </Button>
      {loading ? <Spin /> : (
        bookings.length === 0 ? (
          <Card className="booking-card">暂无预约</Card>
        ) : (
          <List
            dataSource={bookings}
            renderItem={item => {
              const coachName = coaches.find(c => c.id === item.coach_id)?.name || '-';
              const venueName = item.venue_name || venues.find(v => v.id === item.venue_id)?.name || '-';
              return (
                <Card className="booking-card" style={{ marginBottom: 16 }}>
                  <div className="booking-list-row" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <div>
                      <b>日期：</b>{item.date} &nbsp;
                      <b>时间段：</b><Tag color="blue">{item.time_slots}</Tag>
                    </div>
                    <div>
                      <b>场地：</b>{venueName}
                    </div>
                    <div>
                      <b>学员：</b>{item.student_name} {item.student_id ? `(ID: ${item.student_id})` : ''}
                    </div>
                    <div>
                      <b>教练：</b>{coachName}
                    </div>
                    <Button type="link" onClick={() => handleEdit(item)} key="edit">编辑</Button>
                  </div>
                </Card>
              );
            }}
          />
        )
      )}
      <Modal
        title={editing ? '编辑预约' : '新增预约'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleOk}
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="student_name" label="学员姓名" rules={[{ required: true, message: '请输入学员姓名' }]}> 
            <Input className="booking-input" />
          </Form.Item>
          {role === 'admin' && (
            <Form.Item name="coach_id" label="教练" rules={[{ required: true, message: '请选择教练' }]}> 
              <Select className="booking-select" onChange={handleCoachChangeInModal}> 
                {coaches.map(coach => (
                  <Select.Option key={coach.id} value={coach.id}>{coach.name}</Select.Option>
                ))}
              </Select>
            </Form.Item>
          )}
          <Form.Item name="date" label="日期" rules={[{ required: true, message: '请选择日期' }]}> 
            <DatePicker className="booking-picker" onChange={handleDateChangeInModal} />
          </Form.Item>
          <Form.Item name="venue_id" label="场地" rules={[{ required: true, message: '请选择场地' }]}>
            <Select className="booking-select" onChange={handleVenueChangeInModal}>
              {venues.map(venue => (
                <Select.Option key={venue.id} value={venue.id}>{venue.name}</Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item label="时间段" required>
            {(() => {
              const venueId = form.getFieldValue('venue_id');
              if (!venueId) {
                return <div style={{ color: 'red' }}>请先选择场地</div>;
              }
              const date = form.getFieldValue('date');
              if (!date) {
                return <div style={{ color: 'red' }}>请先选择日期</div>;
              }
              const coachId = resolveCoachIdForForm();
              if (!coachId) {
                return <div style={{ color: 'red' }}>请先选择教练</div>;
              }
              return (
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                  {TIME_SLOTS.map(slot => (
                    <Button
                      key={slot}
                      type={selectedSlots.includes(slot) ? 'primary' : 'default'}
                      onClick={() => unavailableSlots.includes(slot) ? null : handleSlotClick(slot)}
                      style={{ marginBottom: 8 }}
                      disabled={unavailableSlots.includes(slot)}
                      className="booking-slot-btn"
                    >
                      {slot}
                    </Button>
                  ))}
                </div>
              );
            })()}
            {selectedSlots.length === 0 && <div style={{ color: 'red', marginTop: 4 }}>请选择至少一个时间段</div>}
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
} 
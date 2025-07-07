import React, { useEffect, useState } from 'react';
import { Card, Button, Select, message, Spin, List, Tag, Modal, Form, Input, DatePicker } from 'antd';
import request from '../utils/request';
import dayjs from 'dayjs';
import { getRole, getUserId } from '../utils/auth';

const TIME_SLOTS = [];
for (let hour = 9; hour < 18; hour++) {
  TIME_SLOTS.push(`${hour.toString().padStart(2, '0')}:00-${hour.toString().padStart(2, '0')}:30`);
  TIME_SLOTS.push(`${hour.toString().padStart(2, '0')}:30-${(hour + 1).toString().padStart(2, '0')}:00`);
}

export default function BookingPage() {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [coaches, setCoaches] = useState([]);
  const [selectedCoach, setSelectedCoach] = useState(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form] = Form.useForm();
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

  // 获取预约列表
  const fetchBookings = async (coachId) => {
    setLoading(true);
    try {
      const res = await request.get('/api/bookings', {
        params: coachId ? { coach_id: coachId } : {},
      });
      setBookings(res.data || []);
    } catch (e) {
      message.error('获取预约列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCoaches();
  }, []);

  useEffect(() => {
    if (selectedCoach) {
      fetchBookings(selectedCoach);
    }
  }, [selectedCoach]);

  // 以日期分组
  const bookingsByDate = bookings.reduce((acc, cur) => {
    const date = cur.date;
    if (!acc[date]) acc[date] = [];
    acc[date].push(cur);
    return acc;
  }, {});

  // 教练选择器
  const coachSelector = role === 'admin' ? (
    <div style={{ marginBottom: 16 }}>
      <span style={{ marginRight: 8 }}>教练：</span>
      <Select
        style={{ width: 200 }}
        value={selectedCoach}
        onChange={setSelectedCoach}
      >
        {coaches.map(coach => (
          <Select.Option key={coach.id} value={coach.id}>{coach.name}</Select.Option>
        ))}
      </Select>
    </div>
  ) : null;

  // 新增/编辑预约弹窗的时间段多选
  const [selectedSlots, setSelectedSlots] = useState([]);
  const [unavailableSlots, setUnavailableSlots] = useState([]);

  // 查询某教练某天已预约的所有时间段
  const fetchUnavailableSlots = async (coachId, date, editingId) => {
    if (!coachId || !date) {
      setUnavailableSlots([]);
      return;
    }
    try {
      const res = await request.get('/api/bookings', {
        params: { coach_id: coachId },
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

  // 新增/编辑预约
  const handleAdd = () => {
    setEditing(null);
    form.resetFields();
    setSelectedSlots([]);
    setModalOpen(true);
    // 拉取当前教练和日期下的已预约时间段
    const coachId = role === 'admin' ? selectedCoach : (coaches.find(c => String(c.user_id) === String(userId))?.id);
    const date = form.getFieldValue('date') ? form.getFieldValue('date').format('YYYY-MM-DD') : dayjs().format('YYYY-MM-DD');
    fetchUnavailableSlots(coachId, date);
  };

  const handleEdit = (record) => {
    setEditing(record);
    form.setFieldsValue({
      ...record,
      date: dayjs(record.date),
    });
    setSelectedSlots(record.time_slots ? record.time_slots.split(',').map(s => s.trim()) : []);
    setModalOpen(true);
    // 拉取当前教练和日期下的已预约时间段（排除自己）
    const coachId = record.coach_id;
    const date = record.date;
    fetchUnavailableSlots(coachId, date, record.id);
  };

  // 日期或教练变更时，重新拉取不可用时间段
  const handleDateOrCoachChange = (dateVal) => {
    const coachId = role === 'admin' ? form.getFieldValue('coach_id') || selectedCoach : (coaches.find(c => String(c.user_id) === String(userId))?.id);
    const date = dateVal ? dateVal.format('YYYY-MM-DD') : form.getFieldValue('date')?.format('YYYY-MM-DD');
    fetchUnavailableSlots(coachId, date, editing?.id);
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
      fetchBookings(selectedCoach);
    } catch (e) {
      message.error('保存失败');
    }
  };

  return (
    <div>
      {coachSelector}
      <Button type="primary" style={{ marginBottom: 16 }} onClick={handleAdd}>
        新增预约
      </Button>
      {loading ? <Spin /> : (
        Object.keys(bookingsByDate).length === 0 ? (
          <Card>暂无预约</Card>
        ) : (
          Object.keys(bookingsByDate).sort().map(date => (
            <Card key={date} title={date} style={{ marginBottom: 16 }}>
              <List
                dataSource={bookingsByDate[date]}
                renderItem={item => (
                  <List.Item
                    actions={[
                      <Button type="link" onClick={() => handleEdit(item)} key="edit">编辑</Button>
                    ]}
                  >
                    <div style={{ width: '100%' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div>
                          <b>学员：</b>{item.student_name} &nbsp;
                          <b>时间段：</b><Tag color="blue">{item.time_slots}</Tag>
                        </div>
                        <div>
                          <b>教练：</b>{coaches.find(c => c.id === item.coach_id)?.name || '-'}
                        </div>
                      </div>
                    </div>
                  </List.Item>
                )}
              />
            </Card>
          ))
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
            <Input />
          </Form.Item>
          {role === 'admin' && (
            <Form.Item name="coach_id" label="教练" rules={[{ required: true, message: '请选择教练' }]}> 
              <Select onChange={() => handleDateOrCoachChange()}> 
                {coaches.map(coach => (
                  <Select.Option key={coach.id} value={coach.id}>{coach.name}</Select.Option>
                ))}
              </Select>
            </Form.Item>
          )}
          <Form.Item name="date" label="日期" rules={[{ required: true, message: '请选择日期' }]}> 
            <DatePicker onChange={handleDateOrCoachChange} />
          </Form.Item>
          <Form.Item label="时间段" required>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
              {TIME_SLOTS.map(slot => (
                <Button
                  key={slot}
                  type={selectedSlots.includes(slot) ? 'primary' : 'default'}
                  onClick={() => unavailableSlots.includes(slot) ? null : handleSlotClick(slot)}
                  style={{ marginBottom: 8 }}
                  disabled={unavailableSlots.includes(slot)}
                >
                  {slot}
                </Button>
              ))}
            </div>
            {selectedSlots.length === 0 && <div style={{ color: 'red', marginTop: 4 }}>请选择至少一个时间段</div>}
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
} 
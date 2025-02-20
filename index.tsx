import { useEffect, useState } from 'react';

export default function Home() {
  const [tasks, setTasks] = useState([]);
  const [newTask, setNewTask] = useState({ title: '', description: '', assignedTo: '' });

  useEffect(() => {
    fetchTasks();
    setupWebSocket();
  }, []);

  const fetchTasks = async () => {
    const res = await fetch('/api/tasks');
    const data = await res.json();
    setTasks(data);
  };

  const setupWebSocket = () => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    ws.onmessage = (event) => {
      setTasks(JSON.parse(event.data));
    };
  };

  const createTask = async () => {
    await fetch('/api/tasks', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newTask),
    });
    setNewTask({ title: '', description: '', assignedTo: '' });
  };

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Task Management System</h1>
      <div className="mb-4">
        <input
          type="text"
          placeholder="Title"
          value={newTask.title}
          onChange={(e) => setNewTask({ ...newTask, title: e.target.value })}
          className="border p-2 mr-2"
        />
        <input
          type="text"
          placeholder="Description"
          value={newTask.description}
          onChange={(e) => setNewTask({ ...newTask, description: e.target.value })}
          className="border p-2 mr-2"
        />
        <input
          type="text"
          placeholder="Assigned To"
          value={newTask.assignedTo}
          onChange={(e) => setNewTask({ ...newTask, assignedTo: e.target.value })}
          className="border p-2 mr-2"
        />
        <button onClick={createTask} className="bg-blue-500 text-white p-2">
          Create Task
        </button>
      </div>
      <div>
        {Object.values(tasks).map((task) => (
          <div key={task.id} className="border p-4 mb-2">
            <h2 className="font-bold">{task.title}</h2>
            <p>{task.description}</p>
            <p>Assigned to: {task.assignedTo}</p>
            <p>Status: {task.status}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

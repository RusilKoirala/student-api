import axios from "axios";

const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:3000";

const api = axios.create({ baseURL: API_BASE });

// Attach JWT from localStorage to every request
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// Auto-logout on 401
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("school");
      window.dispatchEvent(new Event("auth:logout"));
    }
    return Promise.reject(err);
  }
);

// ── Auth 

export const register = (data) => api.post("/api/auth/register", data).then((r) => r.data);
export const login    = (data) => api.post("/api/auth/login",    data).then((r) => r.data);

// ── Students
export const getStudents   = ()     => api.get("/api/students").then((r) => r.data);
export const getStudent    = (id)   => api.get(`/api/students/${id}`).then((r) => r.data);
export const createStudent = (data) => api.post("/api/students", data).then((r) => r.data);
export const updateStudent = (id, data) => api.put(`/api/students/${id}`, data).then((r) => r.data);
export const deleteStudent = (id)   => api.delete(`/api/students/${id}`).then((r) => r.data);

// ── Teachers 
export const getTeachers   = ()     => api.get("/api/teachers").then((r) => r.data);
export const getTeacher    = (id)   => api.get(`/api/teachers/${id}`).then((r) => r.data);
export const createTeacher = (data) => api.post("/api/teachers", data).then((r) => r.data);
export const updateTeacher = (id, data) => api.put(`/api/teachers/${id}`, data).then((r) => r.data);
export const deleteTeacher = (id)   => api.delete(`/api/teachers/${id}`).then((r) => r.data);

// ── Classes 
export const getClasses   = ()     => api.get("/api/classes").then((r) => r.data);
export const getClass     = (id)   => api.get(`/api/classes/${id}`).then((r) => r.data);
export const createClass  = (data) => api.post("/api/classes", data).then((r) => r.data);
export const updateClass  = (id, data) => api.put(`/api/classes/${id}`, data).then((r) => r.data);
export const deleteClass  = (id)   => api.delete(`/api/classes/${id}`).then((r) => r.data);

// ── Stats
export const getStats = () => api.get("/api/stats").then((r) => r.data);

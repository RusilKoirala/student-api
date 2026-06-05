import { useState } from "react";
import axios from "axios";

function App() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [age, setAge] = useState("");
  
  const [selectedStudent, setSelectedStudent] = useState(null);
  const [students, setStudents] = useState([]);
  const [id, setId] = useState("");

  // CREATE
  async function handleSubmit(e) {
    e.preventDefault();

    await axios.post("http://localhost:3000/api/students", {
      name,
      email,
      age: Number(age),
    });

    alert("Student created");
  }

  // GET ALL
  async function getStudents() {
    const res = await axios.get("http://localhost:3000/api/students");
    setStudents(res.data);
  }

  // GET BY ID
  async function getById() {
  try {
    const res = await axios.get(
      `http://localhost:3000/api/students/${id}`
    );

    setSelectedStudent(res.data);
  } catch (err) {
    // if 404 or not found → clear UI
    setSelectedStudent(null);
  }
}

  // DELETE
  async function deleteStudent() {
    await axios.delete(
      `http://localhost:3000/api/students/${id}`
    );

    alert("Deleted");
  }

  return (
    <>
      <h1>Student API Portal</h1>

     
      <h3>Create</h3>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <br />

        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <br />

        <input
          type="number"
          placeholder="Age"
          value={age}
          onChange={(e) => setAge(e.target.value)}
        />
        <br />

        <button type="submit">Create</button>
      </form>

      <hr />
    <div>
      <h3>Student ID</h3>
      <input
        type="text"
        placeholder="Enter ID"
        value={id}
        onChange={(e) => setId(e.target.value)}
      /> <br/>

      <button onClick={getById}>Get By ID</button><br/>
      <button onClick={deleteStudent}>Delete</button>
    <div>
     <h3>Student Details</h3>

{selectedStudent && selectedStudent.id ? (
  <div>
    <p><b>ID:</b> {selectedStudent.id}</p>
    <p><b>Name:</b> {selectedStudent.name}</p>
    <p><b>Email:</b> {selectedStudent.email}</p>
    <p><b>Age:</b> {selectedStudent.age}</p>
  </div>
) : <h3>Not found</h3>} 
    </div>
    </div>
      <hr />

    <div>
      <h3>All Students</h3>
      <button onClick={getStudents}>Load Students</button>
    <div>
      <ul>
        {students.map((s) => (
          <li key={s.id}>
            {s.name} - {s.email} - {s.age}
          </li>
        ))}
      </ul>
    </div>
    </div>
    </>
  );
}

export default App;
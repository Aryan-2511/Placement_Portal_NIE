import axios from 'axios';
const API_URL = 'http://127.0.0.1:8080';
export async function login({ email, password, role }) {
  try {
    const data = { email, password, role };
    const response = await axios.post(`${API_URL}/login`, data, {
      headers: {
        'Content-Type': 'application/json',
      },
    });

    return response.data;
  } catch (err) {
    console.error(err);
    throw err;
  }
}

export async function getStudentDetails(usn) {
  if (!usn) throw new Error('USN is required');
  try {
    const response = await axios.get('http://127.0.0.1:8080/student/details', {
      params: { usn },
      headers: {
        'Content-Type': 'application/json',
      },
    });

    return response.data;
  } catch (error) {
    console.error('Failed to fetch student details:', error);
    if (error.response) {
      throw new Error(`Error ${error.response.status}: ${error.response.data}`);
    } else if (error.request) {
      throw new Error('No response received from the server');
    } else {
      throw new Error(error.message);
    }
  }
}

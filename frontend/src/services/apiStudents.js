import axios from 'axios';
import Cookies from 'js-cookie';

const API_URL = 'http://127.0.0.1:8080/students';

export async function getAllStudentsByBatch(batch) {
  try {
    const user = Cookies.get('user');
    if (!user) throw new Error('User not logged in or initialized');
    const { token } = JSON.parse(user);
    if (!batch) {
      throw new Error('Batch is required');
    }
    if (!token) {
      throw new Error('Unauthorized access not allowed!');
    }
    const response = await axios.get(`${API_URL}/all-by-batch`, {
      params: { batch },
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data; // Return the list of all students by batch
  } catch (error) {
    console.error('Failed to fetch all students by batch:', error);
    if (error.response) {
      throw new Error(`Error ${error.response.status}: ${error.response.data}`);
    } else if (error.request) {
      throw new Error('No response received from the server');
    } else {
      throw new Error(error.message);
    }
  }
}

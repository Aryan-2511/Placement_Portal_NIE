import axios from 'axios';
import Cookies from 'js-cookie';

const API_URL = 'http://127.0.0.1:8080';

export async function addEvent(event) {
  if (!event) throw new Error('Event cannot be empty!');
  try {
    const { token } = JSON.parse(Cookies.get('user'));
    const response = await axios.post(`${API_URL}/schedule/add`, event, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    console.error('Failed to add event:', error);
    if (error.response) {
      throw new Error(`Error ${error.response.status}: ${error.response.data}`);
    } else if (error.request) {
      throw new Error('No response received from the server');
    } else {
      throw new Error(error.message);
    }
  }
}

export async function getAllEvents() {
  try {
    const { token } = JSON.parse(Cookies.get('user'));
    const response = await axios.get(`${API_URL}/schedule/all`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    console.error('Failed to fetch events:', error);
    if (error.response) {
      throw new Error(`Error ${error.response.status}: ${error.response.data}`);
    } else if (error.request) {
      throw new Error('No response received from the server');
    } else {
      throw new Error(error.message);
    }
  }
}

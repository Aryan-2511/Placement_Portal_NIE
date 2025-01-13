import axios from 'axios';
import Cookies from 'js-cookie';

const API_URL = 'http://127.0.0.1:8080/opportunities';
export async function getOpportunitiesByBatch(batch) {
  const user = Cookies.get('user');
  if (!user) throw new Error('User not logged in or initialized');
  const { token } = JSON.parse(user);
  if (!batch) {
    throw new Error('Batch is required');
  }
  if (!token) {
    throw new Error('Unauthorized access not allowed!');
  }

  try {
    const response = await axios.get(`${API_URL}/by-batch`, {
      params: { batch },
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data; // Return the list of opportunities
  } catch (error) {
    console.error('Failed to fetch opportunities by batch:', error);
    if (error.response) {
      throw new Error(`Error ${error.response.status}: ${error.response.data}`);
    } else if (error.request) {
      throw new Error('No response received from the server');
    } else {
      throw new Error(error.message);
    }
  }
}

export async function getOpportunity(id) {
  const user = Cookies.get('user');
  if (!user) throw new Error('User not logged in or initialized');
  const { token } = JSON.parse(user);

  if (!id) throw new Error('Opportunity id is required');
  if (!token) {
    throw new Error('Unauthorized access not allowed!');
  }
  try {
    const response = await axios.get(`${API_URL}/details`, {
      params: { id },
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    console.error('Failed to fetch opportunity:', error);
    if (error.response) {
      throw new Error(`Error ${error.response.status}: ${error.response.data}`);
    } else if (error.request) {
      throw new Error('No response received from the server');
    } else {
      throw new Error(error.message);
    }
  }
}

export async function updateCurrentOpportunity(id, updatedOpportunity) {
  const { role } = JSON.parse(Cookies.get('user'));
  if (!role) throw new Error('Unauthorized access not allowed!');
  try {
    const response = await axios.put(`${API_URL}/edit`, updatedOpportunity, {
      headers: {
        'Content-Type': 'application/json',
        Role: role,
      },
      params: { opportunity_id: id },
    });
    return response.data;
  } catch (error) {
    console.error('Failed to update opportunity:', error);
    if (error.response) {
      throw new Error(`Error ${error.response.status}: ${error.response.data}`);
    } else if (error.request) {
      throw new Error('No response received from the server');
    } else {
      throw new Error(error.message);
    }
  }
}

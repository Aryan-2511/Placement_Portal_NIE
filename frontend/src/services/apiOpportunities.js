import { useRole } from '@/context/UserRoleContext';
import axios from 'axios';

const API_URL = 'http://127.0.0.1:8080/opportunities';
export async function getOpportunitiesByBatch(batch) {
  if (!batch) {
    throw new Error('Batch is required');
  }

  try {
    const response = await axios.get(`${API_URL}/by-batch`, {
      params: { batch },
      headers: {
        'Content-Type': 'application/json',
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
  if (!id) throw new Error('id is required');
  try {
    const response = await axios.get(`${API_URL}/details`, {
      params: { id },
      headers: {
        'Content-Type': 'application/json',
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

export async function updateOpportunity(id, updatedOpportunity, role) {
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

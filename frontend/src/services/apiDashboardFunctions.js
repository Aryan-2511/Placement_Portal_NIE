import axios from 'axios';

const API_URL = 'http://127.0.0.1:8080/student-dash';

export async function getActiveOpportunities(token, batch) {
  try {
    if (!batch) {
      throw new Error(
        'Failed to fetch active opportunities as batch is not provided!'
      );
    }
    if (!token) {
      throw new Error('Unauthorized access not allowed!');
    }
    const response = await axios.get(`${API_URL}/active-opportunities`, {
      params: { batch },
      headers: { Authorization: `Bearer ${token}` },
    });

    return response.data;
  } catch (error) {
    console.error('Failed to fetch active opportunities', error.message);
  }
}
export async function getRecentOpportunities(token, batch) {
  try {
    if (!batch) {
      throw new Error(
        'Failed to fetch recent opportunities as batch is not provided!'
      );
    }
    if (!token) {
      throw new Error('Unauthorized access not allowed!');
    }
    const response = await axios.get(`${API_URL}/recent-opportunities`, {
      params: { batch },
      headers: { Authorization: `Bearer ${token}` },
    });

    return response.data;
  } catch (error) {
    console.error('Failed to fetch recent opportunities', error.message);
  }
}
export async function getTotalPlacedStudents(token, batch) {
  try {
    if (!batch) {
      throw new Error(
        'Failed to fetch total placed students as batch is not provided!'
      );
    }
    if (!token) {
      throw new Error('Unauthorized access not allowed!');
    }
    const response = await axios.get(`${API_URL}/total-placed-students`, {
      params: { batch },
      headers: { Authorization: `Bearer ${token}` },
    });

    return response.data;
  } catch (error) {
    console.error('Failed to total placed students', error.message);
  }
}
export async function getTotalApplications(token, usn) {
  try {
    if (!usn) {
      throw new Error(
        'Failed to fetch total applications as usn is not provided!'
      );
    }
    if (!token) {
      throw new Error('Unauthorized access not allowed!');
    }
    const response = await axios.get(`${API_URL}/total-applications`, {
      params: { usn },
      headers: { Authorization: `Bearer ${token}` },
    });

    return response.data;
  } catch (error) {
    console.error('Failed to total applications', error.message);
  }
}

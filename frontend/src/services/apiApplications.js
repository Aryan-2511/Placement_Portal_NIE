import axios from 'axios';

const API_URL = 'http://127.0.0.1:8080/applications';

export async function getApplications(usn, token) {
  try {
    // const token = getItem('currentUser').token;
    if (!usn) {
      throw new Error('USN is required');
    }
    if (!token) {
      throw new Error('Unauthorized access not allowed!');
    }
    const response = await axios.get(`${API_URL}/student`, {
      params: { usn },
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    console.error('Failed to fetch applications:', error);
    if (error.response) {
      throw new Error(`Error ${error.response.status}: ${error.response.data}`);
    } else if (error.request) {
      throw new Error('No response received from the server');
    } else {
      throw new Error(error.message);
    }
  }
}

import axios from 'axios';

const API_URL = 'http://127.0.0.1:8080/applications';

export async function getApplications(usn) {
  // code
  try {
    const response = await axios.get(`${API_URL}/student`, {
      params: { usn },
      headers: {
        'Content-Type': 'application/json',
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

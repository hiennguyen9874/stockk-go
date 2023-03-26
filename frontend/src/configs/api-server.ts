let API_SERVER_VAL = 'http://localhost:5000';

switch (process.env.NODE_ENV) {
  case 'development':
    // if (process.env.REACT_APP_API_SERVER) {
    //   API_SERVER_VAL = process.env.REACT_APP_API_SERVER;
    // }
    break;
  case 'production':
    // if (process.env.REACT_APP_API_SERVER) {
    //   API_SERVER_VAL = process.env.REACT_APP_API_SERVER;
    // }
    API_SERVER_VAL = '';
    break;
  default:
    break;
}

// export const API_SERVER = API_SERVER_VAL;
export const API_SERVER = API_SERVER_VAL;

export const API_DATAFEED_URL = `${API_SERVER}/api/dchart`;

export const API_STORAGE_URL = `${API_SERVER}/api/storage`;

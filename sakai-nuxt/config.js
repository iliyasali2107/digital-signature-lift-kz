export const API_URL = import.meta.env.VITE_APP_API_URL;
//"http://localhost:5001";
export const getHeader = function (access_token = '') {
    return {
        'Content-Type': 'application/json',
        Accept: '*/*',
        Authorization: `Bearer ${access_token}`
    };
};

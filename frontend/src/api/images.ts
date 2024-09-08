import session from '@/api/session.js';

export default {
  upload(formData: FormData) {
    return session.post('images/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  }
};
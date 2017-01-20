import fetch from 'isomorphic-fetch';

export const fetchUser = () => (
  {
    type: 'FETCH_USER',
  }
);

export const fetchUserFailure = error => (
  {
    type: 'FETCH_USER_FAILURE',
    error,
  }
);

export const fetchUserSuccess = response => (
  {
    type: 'FETCH_USER_SUCCESS',
    response,
  }
);

export const fetchCurrentUser = () => (
  (dispatch) => {
    dispatch(fetchUser());

    return fetch('/manage/api/whoami', {
      credentials: 'include',
    })
      .then(response => Promise.all([response, response.json()]))
      .then(([response, json]) => {
        if (!response.ok) {
          const error = new Error(json.message);
          error.response = response;
          throw error;
        }
        console.log('user json', response, json);
        return json;
      })
      .then(json => (
        dispatch(fetchUserSuccess(json))
      ))
      .catch((err) => {
        console.log('err', err.response);
        if (err.response && err.response.status === 401) {
          window.location = `/manage/api/login?back=${encodeURIComponent(window.location)}`;
        }
        return dispatch(fetchUserFailure(err));
      });
  }
);


export const fetchPlaces = () => (
  {
    type: 'FETCH_PLACES',
  }
);

export const fetchPlacesFailure = error => (
  {
    type: 'FETCH_PLACES_FAILURE',
    error,
  }
);

export const fetchPlacesSuccess = response => (
  {
    type: 'FETCH_PLACES_SUCCESS',
    response,
  }
);

export const fetchAllPlaces = () => (
  (dispatch) => {
    dispatch(fetchUser());

    return fetch('/manage/api/places', {
      credentials: 'include',
    })
      .then(response => response.json())
      .then(json => (
        dispatch(fetchPlacesSuccess(json))
      ))
      .catch(err => (
        dispatch(fetchPlacesFailure(err))
      ));
  }
);

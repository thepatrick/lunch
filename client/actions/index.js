import fetch from 'isomorphic-fetch';
import { push } from 'react-router-redux';

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
        // ('user json', response, json);
        return json;
      })
      .then(json => (
        dispatch(fetchUserSuccess(json))
      ))
      .catch((err) => {
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

const savePlace = id => (
  {
    type: 'SAVE_PLACE',
    id,
  }
);

const savePlaceSuccess = (id, placeName) => (
  {
    type: 'SAVE_PLACE_SUCCESS',
    id,
    placeName,
  }
);

const savePlaceFailure = (id, error) => (
  {
    type: 'SAVE_PLACE_FAILURE',
    id,
    error,
  }
);

export const updatePlaceName = (placeId, placeName) => (
  (dispatch) => {
    dispatch(savePlace(placeId));

    return fetch(`/manage/api/places/${placeId}`, {
      credentials: 'include',
      method: 'POST',
      body: JSON.stringify({
        name: placeName,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(response.statusText);
        }
        return response;
      })
      .then((/* response */) => (
        dispatch(savePlaceSuccess(placeId, placeName))
      ))
      .catch(err => dispatch(savePlaceFailure(placeId, err)));
  }
);

const deletingPlace = id => (
  {
    type: 'DELETE_PLACE',
    id,
  }
);

const deletePlaceSuccess = id => (
  {
    type: 'DELETE_PLACE_SUCCESS',
    id,
  }
);

const deletePlaceFailure = (id, error) => (
  {
    type: 'DELETE_PLACE_FAILURE',
    id,
    error,
  }
);

export const deletePlace = placeId => (
  (dispatch) => {
    dispatch(deletingPlace);

    return fetch(`/manage/api/places/${placeId}`, {
      credentials: 'include',
      method: 'DELETE',
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(response.statusText);
        }
        return response;
      })
      .then(() => {
        dispatch(deletePlaceSuccess(placeId));
        dispatch(push('/manage'));
      })
      .catch(err => dispatch(deletePlaceFailure(placeId, err)));
  }
);

// "user": {
//     "isFetching": false,
//     "error": null,
//     "name": "",
//     "team": {
//         "id": "FOO-BAR",
//         "name": "foobar"
//     }
// },

const defaultState = {
  isFetching: false,
  error: null,
  name: '',
  id: '',
  teamId: '',
  teamName: '',
  teamDomain: '',
};

const user = (state = defaultState, action) => {
  switch (action.type) {
    case 'FETCH_USER':
      return {
        isFetching: true,
        error: null,
        name: '',
        id: '',
        teamId: '',
        teamName: '',
        teamDomain: '',
      };
    case 'FETCH_USER_SUCCESS':
      return {
        isFetching: false,
        error: null,
        name: action.response.user.name,
        id: action.response.user.id,
        teamId: action.response.team.id,
        teamName: action.response.team.name,
        teamDomain: action.response.team.domain,
      };
    case 'FETCH_USER_FAILURE':
      return Object.assign({}, state, {
        isFetching: false,
        error: action.error,
      });
    default:
      return state;
  }
};

export default user;

// const todo = (state = {}, action) => {
//   switch (action.type) {
//     case 'ADD_TODO':
//       return {
//         id: action.id,
//         text: action.text,
//         completed: false,
//       }
//     case 'TOGGLE_TODO':
//       if (state.id !== action.id) {
//         return state;
//       }

//       return Object.assign({}, state, {
//         completed: !state.completed,
//       });

//     default:
//       return state;
//   }
// };

// const todos = (state = [], action) => {
//   switch (action.type) {
//     case 'ADD_TODO':
//       return [
//         ...state,
//         todo(undefined, action),
//       ];
//     case 'TOGGLE_TODO':
//       return state.map(t => todo(t, action));
//     default:
//       return state;
//   }
// };

// export default todos;

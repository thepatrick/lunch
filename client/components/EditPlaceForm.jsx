import React from 'react';
import PropTypes from 'prop-types';
import { Field, reduxForm } from 'redux-form';

const EditPlaceForm = ({ isFetching, isSaving, handleSubmit }) => {
  if (isFetching) {
    return (
      <div>Loading...</div>
    );
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="form-group row">
        <label htmlFor="name-input-field" className="col-2 col-form-label">Name</label>
        <div className="col-10">
          <Field
            className="form-control"
            name="name"
            component="input"
            type="text"
            id="name-input-field"
          />
        </div>
      </div>
      <div className="form-group row">
        <div className="col-10 offset-2">
          <button type="submit" className="btn btn-primary" disabled={isSaving}>Save</button>
        </div>
      </div>
    </form>
  );
};

EditPlaceForm.propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  isFetching: PropTypes.bool.isRequired,
  isSaving: PropTypes.bool.isRequired,
};

export default reduxForm({
  form: 'editPlace',
})(EditPlaceForm);

/* number, shape, string, bool, arrayOf, func */
/* onClick={() => onTodoClick(todo.id)} */

import React, {useEffect, useState} from "react";
import {isEmpty} from "ramda";
import {Button, Form} from "tabler-react";
import {FormattedMessage} from "react-intl";
import {inviteUserToOrganization as InviteUserToOrganizationAction} from "../../actions/organizations";
import {useDispatch} from "react-redux";
import {Autocomplete} from "@material-ui/lab";
import {TextField} from "@material-ui/core";

const styles = {
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  marginBottom: "0.5rem",
}

const CreateTeamButton = ({organizationID, allSeasons}) => {
  const dispatch = useDispatch();
  const [isFormOpen, setFormOpen] = useState(false);
  const [season, setSeason] = useState("");
  const [error, setError] = useState(false);
  const inviteUserToOrganization = (organizationID, name, organizationName) =>
    dispatch(InviteUserToOrganizationAction(organizationID, name, organizationName));

  useEffect(() => {
    if (!isEmpty(name)) {
      setError(false);
    }
  }, [name]);

  const handleChange = event => setName(event.target.value);
  const handleFormOpen = event => {
    event.preventDefault();
    setFormOpen(!isFormOpen);
  };

  const submitOrganizationInvite = async event => {
    event.preventDefault();
    if (isEmpty(name)) {
      setError(true);
    } else {
      await inviteUserToOrganization(organizationID, name, organizationName);
      setFormOpen(false);
    }
  }

  return (
    <div style={styles}>
      <Button
        color="success"
        onClick={handleFormOpen}
        icon={"users"}
        size="sm"
        css={{  display: "flex",
          flexDirection: "column",
          alignItems: "center",
          marginBottom: "0.5rem",
          width: "131px"
      }}
      >
        {"Create new team"}
      </Button>
      {isFormOpen && (
        <form onSubmit={submitOrganizationInvite}>
          <Form.FieldSet css={styles}>
            <Form.Group isRequired label="Season" css={styles}>
              <Autocomplete
                freeSolo
                autoComplete
                autoHighlight
                options={[]}
                css={{width: "213px", backgroundColor: "white", height: "38px"}}
                renderInput={(params) => (
                  <TextField {...params}
                             size="small"
                             variant="outlined"

                  />
                )}
              />
            </Form.Group>
            <Form.Group css={styles}>
              <Button type="submit" color="success" className="ml-auto">
                <FormattedMessage id="CreateTeamButton.send" />
              </Button>
            </Form.Group>
          </Form.FieldSet>
        </form>
      )}
    </div>
  );
}

export default CreateTeamButton;

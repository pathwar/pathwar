import {useIntl} from "react-intl";

//TODO: Lister les organisations de l'utilisateur dans un tableau
//TODO: Créer un boutton permettant de créer une organisation pour l'utilisateur
//TODO: Le button devrait ouvrier un modal permettant de créer une organisation
//TODO: Lister les invitations de l'utilisateur dans un tableau
const OrganizationsPage = () => {
  const intl = useIntl();
  const pageTitleIntl = intl.formatMessage({ id: "Organizations.title" });
}

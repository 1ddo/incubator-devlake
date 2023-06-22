from pydevlake.api import API

class ArgocdAPI(API):
    """_summary_

    Args:
        API (_type_): _description_
    """
    def __init__(self, url: str):
        """_summary_

        Args:
            url (str): _description_
        """
        self.url = url

    def test_connection(self):
        """_summary_

        Returns:
            _type_: _description_
        """
        return self.get(f'{self.url}')

    def applications(self, org: str):
        """_summary_

        Args:
            org (str): _description_

        Returns:
            _type_: _description_
        """
        return self.get(f'{self.url}')


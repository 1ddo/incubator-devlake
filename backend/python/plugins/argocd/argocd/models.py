from pydevlake import Field
from pydevlake.model import ToolModel

class User(ToolModel, table=True):
    id: str = Field(primary_key=True)
    name: str
    email: str
    address: str = Field(source="/contactInfo/address")

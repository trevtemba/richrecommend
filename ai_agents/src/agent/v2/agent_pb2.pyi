from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ProductRequest(_message.Message):
    __slots__ = ("json_input",)
    JSON_INPUT_FIELD_NUMBER: _ClassVar[int]
    json_input: str
    def __init__(self, json_input: _Optional[str] = ...) -> None: ...

class ProductResponse(_message.Message):
    __slots__ = ("products",)
    PRODUCTS_FIELD_NUMBER: _ClassVar[int]
    products: _containers.RepeatedCompositeFieldContainer[ParsedProduct]
    def __init__(self, products: _Optional[_Iterable[_Union[ParsedProduct, _Mapping]]] = ...) -> None: ...

class ParsedProduct(_message.Message):
    __slots__ = ("name", "description", "thumbnail", "ingredients", "retailers")
    NAME_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    THUMBNAIL_FIELD_NUMBER: _ClassVar[int]
    INGREDIENTS_FIELD_NUMBER: _ClassVar[int]
    RETAILERS_FIELD_NUMBER: _ClassVar[int]
    name: str
    description: str
    thumbnail: str
    ingredients: _containers.RepeatedScalarFieldContainer[str]
    retailers: _containers.RepeatedCompositeFieldContainer[Retailer]
    def __init__(self, name: _Optional[str] = ..., description: _Optional[str] = ..., thumbnail: _Optional[str] = ..., ingredients: _Optional[_Iterable[str]] = ..., retailers: _Optional[_Iterable[_Union[Retailer, _Mapping]]] = ...) -> None: ...

class Retailer(_message.Message):
    __slots__ = ("name", "link", "rating", "price", "in_stock")
    NAME_FIELD_NUMBER: _ClassVar[int]
    LINK_FIELD_NUMBER: _ClassVar[int]
    RATING_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    IN_STOCK_FIELD_NUMBER: _ClassVar[int]
    name: str
    link: str
    rating: float
    price: float
    in_stock: bool
    def __init__(self, name: _Optional[str] = ..., link: _Optional[str] = ..., rating: _Optional[float] = ..., price: _Optional[float] = ..., in_stock: bool = ...) -> None: ...

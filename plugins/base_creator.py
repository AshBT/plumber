import py2neo
from . import enhancer

class BaseCreator(enhancer.Enhancer):
    def enhance(self, node):
        if 'phone' in node:
            numbers = node['phone'] if isinstance(node['phone'], list) else [node['phone']]
            relationships = []
            for num in numbers:
                phone_node = self.db.merge_one("Entity", "identifier", num)
                relationships.append(py2neo.Relationship(phone_node, "BY_PHONE", node, number=num, user='auto'))
            self.db.create_unique(*relationships)
